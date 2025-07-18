import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useOrganization } from '../context/OrganizationContext';
import { organizationService, CreateOrganizationRequest } from '../services/organizationService';

const CreateOrganization: React.FC = () => {
  const navigate = useNavigate();
  const { refreshOrganizations } = useOrganization();
  
  const [formData, setFormData] = useState<CreateOrganizationRequest>({
    name: '',
    slug: '',
    domain: '',
    plan_type: 'free',
  });
  
  const [isLoading, setIsLoading] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});

  const handleNameChange = (name: string) => {
    const slug = organizationService.generateSlug(name);
    setFormData({
      ...formData,
      name,
      slug,
    });
    
    // Clear name and slug errors when user types
    if (errors.name || errors.slug) {
      setErrors({
        ...errors,
        name: '',
        slug: '',
      });
    }
  };

  const handleSlugChange = (slug: string) => {
    setFormData({
      ...formData,
      slug: slug.toLowerCase().replace(/[^a-z0-9-]/g, ''),
    });
    
    if (errors.slug) {
      setErrors({
        ...errors,
        slug: '',
      });
    }
  };

  const validateForm = (): boolean => {
    const newErrors: Record<string, string> = {};

    if (!formData.name.trim()) {
      newErrors.name = 'Organization name is required';
    } else if (formData.name.trim().length < 2) {
      newErrors.name = 'Organization name must be at least 2 characters';
    } else if (formData.name.trim().length > 255) {
      newErrors.name = 'Organization name must be less than 255 characters';
    }

    if (!formData.slug.trim()) {
      newErrors.slug = 'Organization slug is required';
    } else if (!organizationService.isValidSlug(formData.slug)) {
      newErrors.slug = 'Slug must be 2-100 characters, contain only lowercase letters, numbers, and hyphens, and cannot start or end with a hyphen';
    }

    if (formData.domain && !isValidDomain(formData.domain)) {
      newErrors.domain = 'Please enter a valid domain name';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const isValidDomain = (domain: string): boolean => {
    const domainRegex = /^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\.[a-zA-Z]{2,}$/;
    return domainRegex.test(domain);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    setIsLoading(true);
    try {
      const response = await organizationService.createOrganization(formData);
      const newOrganization = response.data;
      
      // Refresh organizations list
      await refreshOrganizations();
      
      // Switch to the new organization
      await organizationService.switchOrganization(newOrganization.id);
      
      // Navigate to organization management
      navigate('/organizations/manage');
    } catch (error: any) {
      console.error('Failed to create organization:', error);
      
      // Handle specific API errors
      if (error.response?.data?.message) {
        if (error.response.data.message.includes('slug')) {
          setErrors({ slug: 'This slug is already taken. Please choose a different one.' });
        } else {
          setErrors({ general: error.response.data.message });
        }
      } else {
        setErrors({ general: 'Failed to create organization. Please try again.' });
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
      <div className="sm:mx-auto sm:w-full sm:max-w-md">
        <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
          Create New Organization
        </h2>
        <p className="mt-2 text-center text-sm text-gray-600">
          Set up your organization to get started with team collaboration
        </p>
      </div>

      <div className="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
        <div className="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10">
          <form className="space-y-6" onSubmit={handleSubmit}>
            {/* General Error */}
            {errors.general && (
              <div className="rounded-md bg-red-50 p-4">
                <div className="text-sm text-red-700">{errors.general}</div>
              </div>
            )}

            {/* Organization Name */}
            <div>
              <label htmlFor="name" className="block text-sm font-medium text-gray-700">
                Organization Name *
              </label>
              <div className="mt-1">
                <input
                  id="name"
                  name="name"
                  type="text"
                  required
                  value={formData.name}
                  onChange={(e) => handleNameChange(e.target.value)}
                  className={`appearance-none block w-full px-3 py-2 border rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm ${
                    errors.name ? 'border-red-300' : 'border-gray-300'
                  }`}
                  placeholder="Enter organization name"
                />
              </div>
              {errors.name && (
                <p className="mt-2 text-sm text-red-600">{errors.name}</p>
              )}
            </div>

            {/* Organization Slug */}
            <div>
              <label htmlFor="slug" className="block text-sm font-medium text-gray-700">
                Organization Slug *
              </label>
              <div className="mt-1">
                <div className="flex rounded-md shadow-sm">
                  <span className="inline-flex items-center px-3 rounded-l-md border border-r-0 border-gray-300 bg-gray-50 text-gray-500 text-sm">
                    yourapp.com/
                  </span>
                  <input
                    id="slug"
                    name="slug"
                    type="text"
                    required
                    value={formData.slug}
                    onChange={(e) => handleSlugChange(e.target.value)}
                    className={`flex-1 block w-full px-3 py-2 border rounded-none rounded-r-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm ${
                      errors.slug ? 'border-red-300' : 'border-gray-300'
                    }`}
                    placeholder="organization-slug"
                  />
                </div>
              </div>
              {errors.slug && (
                <p className="mt-2 text-sm text-red-600">{errors.slug}</p>
              )}
              <p className="mt-2 text-sm text-gray-500">
                Used in URLs and must be unique. Only lowercase letters, numbers, and hyphens allowed.
              </p>
            </div>

            {/* Domain (Optional) */}
            <div>
              <label htmlFor="domain" className="block text-sm font-medium text-gray-700">
                Domain (Optional)
              </label>
              <div className="mt-1">
                <input
                  id="domain"
                  name="domain"
                  type="text"
                  value={formData.domain}
                  onChange={(e) => setFormData({ ...formData, domain: e.target.value })}
                  className={`appearance-none block w-full px-3 py-2 border rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm ${
                    errors.domain ? 'border-red-300' : 'border-gray-300'
                  }`}
                  placeholder="example.com"
                />
              </div>
              {errors.domain && (
                <p className="mt-2 text-sm text-red-600">{errors.domain}</p>
              )}
              <p className="mt-2 text-sm text-gray-500">
                Your organization's website domain for branding purposes.
              </p>
            </div>

            {/* Plan Type */}
            <div>
              <label htmlFor="plan_type" className="block text-sm font-medium text-gray-700">
                Plan Type
              </label>
              <div className="mt-1">
                <select
                  id="plan_type"
                  name="plan_type"
                  value={formData.plan_type}
                  onChange={(e) => setFormData({ ...formData, plan_type: e.target.value as any })}
                  className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                >
                  <option value="free">Free</option>
                  <option value="basic">Basic</option>
                  <option value="premium">Premium</option>
                  <option value="enterprise">Enterprise</option>
                </select>
              </div>
              <p className="mt-2 text-sm text-gray-500">
                You can upgrade your plan later from organization settings.
              </p>
            </div>

            {/* Submit Button */}
            <div>
              <button
                type="submit"
                disabled={isLoading}
                className={`w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 ${
                  isLoading
                    ? 'bg-gray-400 cursor-not-allowed'
                    : 'bg-indigo-600 hover:bg-indigo-700'
                }`}
              >
                {isLoading ? (
                  <div className="flex items-center">
                    <div className="animate-spin -ml-1 mr-3 h-5 w-5 border-2 border-white border-t-transparent rounded-full"></div>
                    Creating Organization...
                  </div>
                ) : (
                  'Create Organization'
                )}
              </button>
            </div>

            {/* Cancel Link */}
            <div className="text-center">
              <button
                type="button"
                onClick={() => navigate(-1)}
                className="text-sm text-indigo-600 hover:text-indigo-500"
              >
                Cancel and go back
              </button>
            </div>
          </form>
        </div>
      </div>

      {/* Preview Card */}
      {(formData.name || formData.slug) && (
        <div className="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
          <div className="bg-white py-6 px-4 shadow sm:rounded-lg sm:px-6">
            <h3 className="text-lg font-medium text-gray-900 mb-4">Preview</h3>
            <div className="space-y-3">
              {formData.name && (
                <div>
                  <span className="text-sm font-medium text-gray-500">Organization Name:</span>
                  <p className="text-sm text-gray-900">{formData.name}</p>
                </div>
              )}
              {formData.slug && (
                <div>
                  <span className="text-sm font-medium text-gray-500">URL:</span>
                  <p className="text-sm text-gray-900">yourapp.com/{formData.slug}</p>
                </div>
              )}
              {formData.domain && (
                <div>
                  <span className="text-sm font-medium text-gray-500">Domain:</span>
                  <p className="text-sm text-gray-900">{formData.domain}</p>
                </div>
              )}
              <div>
                <span className="text-sm font-medium text-gray-500">Plan:</span>
                <span className={`ml-2 inline-flex px-2 py-1 text-xs font-semibold rounded-full ${organizationService.getPlanColor(formData.plan_type)}`}>
                  {organizationService.formatPlanType(formData.plan_type)}
                </span>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default CreateOrganization;