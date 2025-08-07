package ai

import (
	"fmt"
	"math"
	"sort"
)

// WeightedAverageVoting implements weighted average voting strategy
type WeightedAverageVoting struct{}

// MajorityVoting implements majority voting strategy
type MajorityVoting struct{}

// RankedChoiceVoting implements ranked choice voting strategy
type RankedChoiceVoting struct{}

// AdaptiveVoting implements adaptive voting based on model confidence
type AdaptiveVoting struct{}

// NewWeightedAverageVoting creates a new weighted average voting strategy
func NewWeightedAverageVoting() *WeightedAverageVoting {
	return &WeightedAverageVoting{}
}

// NewMajorityVoting creates a new majority voting strategy
func NewMajorityVoting() *MajorityVoting {
	return &MajorityVoting{}
}

// NewRankedChoiceVoting creates a new ranked choice voting strategy
func NewRankedChoiceVoting() *RankedChoiceVoting {
	return &RankedChoiceVoting{}
}

// NewAdaptiveVoting creates a new adaptive voting strategy
func NewAdaptiveVoting() *AdaptiveVoting {
	return &AdaptiveVoting{}
}

// CombinePredictions combines predictions using weighted average
func (w *WeightedAverageVoting) CombinePredictions(predictions []ModelPrediction, weights map[string]float64) (*EnsemblePrediction, error) {
	if len(predictions) == 0 {
		return nil, fmt.Errorf("no predictions to combine")
	}

	// Separate numeric and categorical predictions
	numericPredictions := make([]ModelPrediction, 0)
	categoricalPredictions := make([]ModelPrediction, 0)

	for _, pred := range predictions {
		if _, ok := pred.Prediction.(float64); ok {
			numericPredictions = append(numericPredictions, pred)
		} else {
			categoricalPredictions = append(categoricalPredictions, pred)
		}
	}

	var finalPrediction interface{}
	var confidence float64
	modelVotes := make(map[string]ModelPrediction)

	// Store all model votes
	for _, pred := range predictions {
		modelVotes[pred.ModelID] = pred
	}

	if len(numericPredictions) > 0 {
		// Handle numeric predictions with weighted average
		weightedSum := 0.0
		totalWeight := 0.0
		confidenceSum := 0.0

		for _, pred := range numericPredictions {
			weight := weights[pred.ModelID]
			if weight == 0 {
				weight = 1.0 / float64(len(numericPredictions)) // Equal weight if not specified
			}

			value := pred.Prediction.(float64)
			weightedSum += value * weight * pred.Confidence
			totalWeight += weight * pred.Confidence
			confidenceSum += pred.Confidence * weight
		}

		if totalWeight > 0 {
			finalPrediction = weightedSum / totalWeight
			confidence = confidenceSum / float64(len(numericPredictions))
		}
	} else if len(categoricalPredictions) > 0 {
		// Handle categorical predictions with weighted voting
		votes := make(map[interface{}]float64)
		totalWeight := 0.0

		for _, pred := range categoricalPredictions {
			weight := weights[pred.ModelID]
			if weight == 0 {
				weight = 1.0 / float64(len(categoricalPredictions))
			}

			weightedVote := weight * pred.Confidence
			votes[pred.Prediction] += weightedVote
			totalWeight += weightedVote
		}

		// Find the prediction with highest weighted vote
		maxVote := 0.0
		for prediction, vote := range votes {
			if vote > maxVote {
				maxVote = vote
				finalPrediction = prediction
			}
		}

		confidence = maxVote / totalWeight
	}

	return &EnsemblePrediction{
		FinalPrediction: finalPrediction,
		Confidence:      confidence,
		ModelVotes:      modelVotes,
		Metadata: map[string]interface{}{
			"strategy": "weighted_average",
			"numeric_models": len(numericPredictions),
			"categorical_models": len(categoricalPredictions),
		},
	}, nil
}

// CombinePredictions combines predictions using majority voting
func (m *MajorityVoting) CombinePredictions(predictions []ModelPrediction, weights map[string]float64) (*EnsemblePrediction, error) {
	if len(predictions) == 0 {
		return nil, fmt.Errorf("no predictions to combine")
	}

	votes := make(map[interface{}]int)
	confidenceSum := make(map[interface{}]float64)
	modelVotes := make(map[string]ModelPrediction)

	// Count votes for each prediction
	for _, pred := range predictions {
		votes[pred.Prediction]++
		confidenceSum[pred.Prediction] += pred.Confidence
		modelVotes[pred.ModelID] = pred
	}

	// Find majority prediction
	maxVotes := 0
	var finalPrediction interface{}
	var confidence float64

	for prediction, voteCount := range votes {
		if voteCount > maxVotes {
			maxVotes = voteCount
			finalPrediction = prediction
			confidence = confidenceSum[prediction] / float64(voteCount)
		}
	}

	// Calculate consensus (percentage of models that agreed)
	consensus := float64(maxVotes) / float64(len(predictions))

	return &EnsemblePrediction{
		FinalPrediction: finalPrediction,
		Confidence:      confidence * consensus, // Adjust confidence by consensus
		ModelVotes:      modelVotes,
		Metadata: map[string]interface{}{
			"strategy": "majority_voting",
			"votes": maxVotes,
			"total_models": len(predictions),
			"consensus": consensus,
		},
	}, nil
}

// CombinePredictions combines predictions using ranked choice voting
func (r *RankedChoiceVoting) CombinePredictions(predictions []ModelPrediction, weights map[string]float64) (*EnsemblePrediction, error) {
	if len(predictions) == 0 {
		return nil, fmt.Errorf("no predictions to combine")
	}

	// Sort predictions by confidence (highest first)
	sortedPredictions := make([]ModelPrediction, len(predictions))
	copy(sortedPredictions, predictions)
	
	sort.Slice(sortedPredictions, func(i, j int) bool {
		return sortedPredictions[i].Confidence > sortedPredictions[j].Confidence
	})

	modelVotes := make(map[string]ModelPrediction)
	for _, pred := range predictions {
		modelVotes[pred.ModelID] = pred
	}

	// Use ranked choice algorithm
	rounds := make([]map[interface{}]float64, 0)
	candidates := make(map[interface{}]bool)
	
	// Initialize candidates
	for _, pred := range sortedPredictions {
		candidates[pred.Prediction] = true
	}

	// Perform elimination rounds
	for len(candidates) > 1 {
		roundVotes := make(map[interface{}]float64)
		
		// Count votes for remaining candidates
		for _, pred := range sortedPredictions {
			if candidates[pred.Prediction] {
				weight := weights[pred.ModelID]
				if weight == 0 {
					weight = 1.0
				}
				roundVotes[pred.Prediction] += weight * pred.Confidence
			}
		}
		
		rounds = append(rounds, roundVotes)
		
		// Find candidate with lowest votes
		minVotes := math.Inf(1)
		var eliminateCandidate interface{}
		
		for candidate, votes := range roundVotes {
			if votes < minVotes {
				minVotes = votes
				eliminateCandidate = candidate
			}
		}
		
		// Eliminate candidate
		delete(candidates, eliminateCandidate)
	}

	// Winner is the remaining candidate
	var finalPrediction interface{}
	var confidence float64
	
	for candidate := range candidates {
		finalPrediction = candidate
		
		// Calculate average confidence for this prediction
		totalConfidence := 0.0
		count := 0
		for _, pred := range predictions {
			if pred.Prediction == candidate {
				totalConfidence += pred.Confidence
				count++
			}
		}
		if count > 0 {
			confidence = totalConfidence / float64(count)
		}
		break
	}

	return &EnsemblePrediction{
		FinalPrediction: finalPrediction,
		Confidence:      confidence,
		ModelVotes:      modelVotes,
		Metadata: map[string]interface{}{
			"strategy": "ranked_choice",
			"rounds": len(rounds),
			"elimination_rounds": rounds,
		},
	}, nil
}

// CombinePredictions combines predictions using adaptive voting
func (a *AdaptiveVoting) CombinePredictions(predictions []ModelPrediction, weights map[string]float64) (*EnsemblePrediction, error) {
	if len(predictions) == 0 {
		return nil, fmt.Errorf("no predictions to combine")
	}

	modelVotes := make(map[string]ModelPrediction)
	for _, pred := range predictions {
		modelVotes[pred.ModelID] = pred
	}

	// Calculate adaptive weights based on confidence and historical performance
	adaptiveWeights := make(map[string]float64)
	totalAdaptiveWeight := 0.0

	for _, pred := range predictions {
		baseWeight := weights[pred.ModelID]
		if baseWeight == 0 {
			baseWeight = 1.0 / float64(len(predictions))
		}

		// Boost weight based on confidence
		confidenceBoost := math.Pow(pred.Confidence, 2) // Square to emphasize high confidence
		
		// Calculate adaptive weight
		adaptiveWeight := baseWeight * confidenceBoost
		adaptiveWeights[pred.ModelID] = adaptiveWeight
		totalAdaptiveWeight += adaptiveWeight
	}

	// Normalize adaptive weights
	for modelID := range adaptiveWeights {
		adaptiveWeights[modelID] /= totalAdaptiveWeight
	}

	// Determine if we should use numeric or categorical combination
	numericPredictions := make([]ModelPrediction, 0)
	for _, pred := range predictions {
		if _, ok := pred.Prediction.(float64); ok {
			numericPredictions = append(numericPredictions, pred)
		}
	}

	var finalPrediction interface{}
	var confidence float64

	if len(numericPredictions) > len(predictions)/2 {
		// Use weighted average for numeric predictions
		weightedSum := 0.0
		totalWeight := 0.0
		confidenceSum := 0.0

		for _, pred := range numericPredictions {
			weight := adaptiveWeights[pred.ModelID]
			value := pred.Prediction.(float64)
			
			weightedSum += value * weight
			totalWeight += weight
			confidenceSum += pred.Confidence * weight
		}

		if totalWeight > 0 {
			finalPrediction = weightedSum / totalWeight
			confidence = confidenceSum / totalWeight
		}
	} else {
		// Use confidence-weighted voting for categorical predictions
		votes := make(map[interface{}]float64)
		totalWeight := 0.0

		for _, pred := range predictions {
			weight := adaptiveWeights[pred.ModelID]
			votes[pred.Prediction] += weight
			totalWeight += weight
		}

		// Find prediction with highest weighted vote
		maxVote := 0.0
		for prediction, vote := range votes {
			if vote > maxVote {
				maxVote = vote
				finalPrediction = prediction
			}
		}

		confidence = maxVote / totalWeight
	}

	// Calculate uncertainty based on prediction spread
	uncertainty := a.calculateUncertainty(predictions, adaptiveWeights)
	adjustedConfidence := confidence * (1.0 - uncertainty)

	return &EnsemblePrediction{
		FinalPrediction: finalPrediction,
		Confidence:      adjustedConfidence,
		ModelVotes:      modelVotes,
		Weights:         adaptiveWeights,
		Metadata: map[string]interface{}{
			"strategy": "adaptive",
			"uncertainty": uncertainty,
			"original_confidence": confidence,
			"adaptive_weights": adaptiveWeights,
		},
	}, nil
}

// calculateUncertainty calculates prediction uncertainty based on model disagreement
func (a *AdaptiveVoting) calculateUncertainty(predictions []ModelPrediction, weights map[string]float64) float64 {
	if len(predictions) < 2 {
		return 0.0
	}

	// For numeric predictions, calculate weighted variance
	numericPredictions := make([]float64, 0)
	numericWeights := make([]float64, 0)

	for _, pred := range predictions {
		if value, ok := pred.Prediction.(float64); ok {
			numericPredictions = append(numericPredictions, value)
			numericWeights = append(numericWeights, weights[pred.ModelID])
		}
	}

	if len(numericPredictions) >= 2 {
		// Calculate weighted mean
		weightedSum := 0.0
		totalWeight := 0.0
		for i, value := range numericPredictions {
			weight := numericWeights[i]
			weightedSum += value * weight
			totalWeight += weight
		}
		weightedMean := weightedSum / totalWeight

		// Calculate weighted variance
		weightedVariance := 0.0
		for i, value := range numericPredictions {
			weight := numericWeights[i]
			weightedVariance += weight * math.Pow(value-weightedMean, 2)
		}
		weightedVariance /= totalWeight

		// Normalize variance to [0, 1] range
		// This is a heuristic - in practice, you'd calibrate based on your data
		normalizedVariance := math.Min(1.0, weightedVariance/math.Pow(weightedMean, 2))
		return normalizedVariance
	}

	// For categorical predictions, calculate disagreement rate
	votes := make(map[interface{}]float64)
	totalWeight := 0.0

	for _, pred := range predictions {
		weight := weights[pred.ModelID]
		votes[pred.Prediction] += weight
		totalWeight += weight
	}

	// Find maximum vote share
	maxVoteShare := 0.0
	for _, vote := range votes {
		voteShare := vote / totalWeight
		if voteShare > maxVoteShare {
			maxVoteShare = voteShare
		}
	}

	// Uncertainty is 1 - max vote share
	return 1.0 - maxVoteShare
}

// GetVotingStrategy returns a voting strategy by name
func GetVotingStrategy(strategyName string) VotingStrategy {
	switch strategyName {
	case "weighted_average":
		return NewWeightedAverageVoting()
	case "majority":
		return NewMajorityVoting()
	case "ranked_choice":
		return NewRankedChoiceVoting()
	case "adaptive":
		return NewAdaptiveVoting()
	default:
		return NewWeightedAverageVoting() // Default strategy
	}
}
