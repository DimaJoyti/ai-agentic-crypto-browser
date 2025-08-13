import React, { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { CheckCircle, Star, TrendingUp, Shield, Zap, Brain } from 'lucide-react';

const BetaSignupPage = () => {
  const [email, setEmail] = useState('');
  const [name, setName] = useState('');
  const [selectedTier, setSelectedTier] = useState('professional');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitted, setSubmitted] = useState(false);

  const betaTiers = [
    {
      id: 'starter',
      name: 'Starter Beta',
      originalPrice: 49,
      betaPrice: 25,
      savings: '50% OFF',
      features: [
        'Basic AI Trading (85%+ accuracy)',
        '3 Trading Strategies',
        'Single Chain Support',
        'Basic Analytics Dashboard',
        'Email Support',
        'Beta Access to New Features'
      ],
      popular: false,
      color: 'border-gray-200'
    },
    {
      id: 'professional',
      name: 'Professional Beta',
      originalPrice: 199,
      betaPrice: 99,
      savings: '50% OFF',
      features: [
        'Advanced AI Trading (85%+ accuracy)',
        '10+ Trading Strategies',
        'Multi-Chain Support (7+ chains)',
        'Advanced Analytics & Predictions',
        'DeFi Integration',
        'Voice Commands (Beta)',
        'Priority Support',
        'Performance Fee Sharing'
      ],
      popular: true,
      color: 'border-blue-500'
    },
    {
      id: 'enterprise',
      name: 'Enterprise Beta',
      originalPrice: 999,
      betaPrice: 499,
      savings: '50% OFF',
      features: [
        'Full Platform Access',
        'Unlimited Strategies',
        'All Chains Supported',
        'Custom AI Models',
        'White-Label Solution (Beta)',
        'Dedicated Support',
        'Custom Integrations',
        'Revenue Sharing Program'
      ],
      popular: false,
      color: 'border-purple-500'
    }
  ];

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);

    try {
      // API call to register beta user
      const response = await fetch('/api/beta/signup', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email,
          name,
          tier: selectedTier,
        }),
      });

      if (response.ok) {
        setSubmitted(true);
      } else {
        throw new Error('Failed to submit');
      }
    } catch (error) {
      console.error('Error submitting beta signup:', error);
      alert('Error submitting signup. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  if (submitted) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 to-purple-50 flex items-center justify-center p-4">
        <Card className="max-w-md w-full text-center">
          <CardHeader>
            <div className="mx-auto w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mb-4">
              <CheckCircle className="w-8 h-8 text-green-600" />
            </div>
            <CardTitle className="text-2xl">Welcome to the Beta!</CardTitle>
            <CardDescription>
              Thank you for joining our exclusive beta program. You'll receive setup instructions within 24 hours.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="bg-blue-50 p-4 rounded-lg">
                <h3 className="font-semibold text-blue-900">What's Next?</h3>
                <ul className="text-sm text-blue-700 mt-2 space-y-1">
                  <li>• Check your email for beta access credentials</li>
                  <li>• Join our exclusive Discord community</li>
                  <li>• Schedule your onboarding call</li>
                  <li>• Start trading with AI in 48 hours</li>
                </ul>
              </div>
              <Button 
                onClick={() => window.location.href = '/dashboard'}
                className="w-full"
              >
                Go to Dashboard
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-purple-50">
      {/* Hero Section */}
      <div className="container mx-auto px-4 py-16">
        <div className="text-center mb-16">
          <Badge className="mb-4 bg-red-100 text-red-800 border-red-200">
            🔥 LIMITED BETA - 50% OFF
          </Badge>
          <h1 className="text-5xl font-bold mb-6 bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
            AI-Powered Crypto Trading
          </h1>
          <p className="text-xl text-gray-600 mb-8 max-w-3xl mx-auto">
            Join the exclusive beta of the world's most advanced AI cryptocurrency trading platform. 
            <strong> 85%+ prediction accuracy</strong> with institutional-grade features.
          </p>
          
          {/* Key Stats */}
          <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-12">
            <div className="bg-white p-6 rounded-lg shadow-sm">
              <Brain className="w-8 h-8 text-blue-600 mx-auto mb-2" />
              <div className="text-2xl font-bold text-gray-900">85%+</div>
              <div className="text-sm text-gray-600">AI Accuracy</div>
            </div>
            <div className="bg-white p-6 rounded-lg shadow-sm">
              <Zap className="w-8 h-8 text-yellow-600 mx-auto mb-2" />
              <div className="text-2xl font-bold text-gray-900">&lt;100ms</div>
              <div className="text-sm text-gray-600">Execution Speed</div>
            </div>
            <div className="bg-white p-6 rounded-lg shadow-sm">
              <TrendingUp className="w-8 h-8 text-green-600 mx-auto mb-2" />
              <div className="text-2xl font-bold text-gray-900">7+</div>
              <div className="text-sm text-gray-600">Blockchains</div>
            </div>
            <div className="bg-white p-6 rounded-lg shadow-sm">
              <Shield className="w-8 h-8 text-purple-600 mx-auto mb-2" />
              <div className="text-2xl font-bold text-gray-900">100%</div>
              <div className="text-sm text-gray-600">Secure</div>
            </div>
          </div>
        </div>

        {/* Pricing Tiers */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mb-16">
          {betaTiers.map((tier) => (
            <Card 
              key={tier.id}
              className={`relative cursor-pointer transition-all hover:shadow-lg ${
                selectedTier === tier.id ? `ring-2 ring-blue-500 ${tier.color}` : tier.color
              } ${tier.popular ? 'scale-105' : ''}`}
              onClick={() => setSelectedTier(tier.id)}
            >
              {tier.popular && (
                <div className="absolute -top-3 left-1/2 transform -translate-x-1/2">
                  <Badge className="bg-blue-600 text-white">
                    <Star className="w-3 h-3 mr-1" />
                    MOST POPULAR
                  </Badge>
                </div>
              )}
              
              <CardHeader className="text-center">
                <Badge className="mb-2 bg-red-100 text-red-800 border-red-200 w-fit mx-auto">
                  {tier.savings}
                </Badge>
                <CardTitle className="text-xl">{tier.name}</CardTitle>
                <div className="space-y-1">
                  <div className="text-3xl font-bold">
                    ${tier.betaPrice}
                    <span className="text-lg text-gray-500 line-through ml-2">
                      ${tier.originalPrice}
                    </span>
                  </div>
                  <div className="text-sm text-gray-600">/month during beta</div>
                </div>
              </CardHeader>
              
              <CardContent>
                <ul className="space-y-3">
                  {tier.features.map((feature, index) => (
                    <li key={index} className="flex items-start">
                      <CheckCircle className="w-4 h-4 text-green-500 mr-2 mt-0.5 flex-shrink-0" />
                      <span className="text-sm">{feature}</span>
                    </li>
                  ))}
                </ul>
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Signup Form */}
        <Card className="max-w-md mx-auto">
          <CardHeader className="text-center">
            <CardTitle>Join the Beta Program</CardTitle>
            <CardDescription>
              Limited spots available. Start trading with AI in 48 hours.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <Input
                  type="text"
                  placeholder="Full Name"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  required
                />
              </div>
              <div>
                <Input
                  type="email"
                  placeholder="Email Address"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                />
              </div>
              <div className="text-center">
                <p className="text-sm text-gray-600 mb-4">
                  Selected: <strong>{betaTiers.find(t => t.id === selectedTier)?.name}</strong>
                  <br />
                  Beta Price: <strong>${betaTiers.find(t => t.id === selectedTier)?.betaPrice}/month</strong>
                </p>
              </div>
              <Button 
                type="submit" 
                className="w-full bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700"
                disabled={isSubmitting}
              >
                {isSubmitting ? 'Joining Beta...' : 'Join Beta Program'}
              </Button>
              <p className="text-xs text-gray-500 text-center">
                Cancel anytime. No long-term commitment. Beta pricing locked for 6 months.
              </p>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default BetaSignupPage;
