import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useAuth } from '../hooks/useAuth';
import { userService } from '../services/api';
import Icon from '../components/common/Icons';
import TabLayout from '../components/common/TabLayout';

const Profile = () => {
  const { t } = useTranslation();
  const { user, updateUser } = useAuth();
  const [activeTab, setActiveTab] = useState('profile');
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');

  // Profile form state
  const [profileData, setProfileData] = useState({
    username: '',
    first_name: '',
    last_name: '',
    email: '',
    bio: '',
    location: '',
    website: '',
    phone: '',
    language: '',
    timezone: '',
    avatar_url: '',
  });

  // Password form state
  const [passwordData, setPasswordData] = useState({
    current_password: '',
    new_password: '',
    confirm_password: '',
  });

  useEffect(() => {
    // Load user profile data
    const loadProfile = async () => {
      try {
        const response = await userService.getProfile();
        const userData = response.data;
        setProfileData({
          username: userData.username || '',
          first_name: userData.first_name || '',
          last_name: userData.last_name || '',
          email: userData.email || '',
          bio: userData.bio || '',
          location: userData.location || '',
          website: userData.website || '',
          phone: userData.phone || '',
          language: userData.language || 'en',
          timezone: userData.timezone || 'UTC',
          avatar_url: userData.avatar_url || '',
        });
      } catch (err) {
        setError(t('profile.errors.loadFailed'));
      }
    };

    loadProfile();
  }, []);

  const handleProfileChange = (e) => {
    setProfileData({
      ...profileData,
      [e.target.name]: e.target.value,
    });
  };

  const handlePasswordChange = (e) => {
    setPasswordData({
      ...passwordData,
      [e.target.name]: e.target.value,
    });
  };

  const handleProfileSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setMessage('');

    try {
      const response = await userService.updateProfile(profileData);
      updateUser(response.data);
      setMessage(t('profile.profileUpdated'));
    } catch (err) {
      setError(err.response?.data?.message || t('profile.errors.updateFailed'));
    } finally {
      setLoading(false);
    }
  };

  const handlePasswordSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setMessage('');

    if (passwordData.new_password !== passwordData.confirm_password) {
      setError(t('profile.errors.passwordMismatch'));
      setLoading(false);
      return;
    }

    if (passwordData.new_password.length < 8) {
      setError(t('profile.errors.passwordMinLength'));
      setLoading(false);
      return;
    }

    try {
      await userService.changePassword(passwordData);
      setMessage(t('profile.passwordChanged'));
      setPasswordData({
        current_password: '',
        new_password: '',
        confirm_password: '',
      });
    } catch (err) {
      setError(err.response?.data?.message || t('profile.errors.changeFailed'));
    } finally {
      setLoading(false);
    }
  };

  const tabs = [
    { id: 'profile', name: t('profile.profileInfoTab'), icon: 'user' },
    { id: 'password', name: t('profile.changePasswordTab'), icon: 'shield' },
  ];

  return (
    <TabLayout
      title={t('profile.userProfile')}
      subtitle={t('profile.profileDescription')}
      showTabs={true}
      tabs={tabs}
      activeTab={activeTab}
      onTabChange={setActiveTab}
      message={message}
      error={error}
      onMessageDismiss={() => setMessage('')}
      onErrorDismiss={() => setError('')}
    >

            {/* Profile Information Tab */}
            {activeTab === 'profile' && (
              <form onSubmit={handleProfileSubmit} className="space-y-6">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <label htmlFor="username" className="block text-sm font-medium text-gray-700 mb-2">
                      {t('profile.username')}
                    </label>
                    <input
                      type="text"
                      id="username"
                      name="username"
                      value={profileData.username}
                      onChange={handleProfileChange}
                      className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors"
                      placeholder={t('profile.usernamePlaceholder')}
                      required
                      minLength={3}
                      maxLength={30}
                      pattern="[a-zA-Z0-9_]+"
                      title={t('profile.usernameHelp')}
                    />
                  </div>

                  <div>
                    <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-2">
                      {t('common.email')}
                    </label>
                    <input
                      type="email"
                      id="email"
                      name="email"
                      value={profileData.email}
                      onChange={handleProfileChange}
                      className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors"
                      placeholder={t('profile.emailPlaceholder')}
                      required
                    />
                  </div>

                  <div>
                    <label htmlFor="first_name" className="block text-sm font-medium text-gray-700 mb-2">
                      {t('users.firstName')}
                    </label>
                    <input
                      type="text"
                      id="first_name"
                      name="first_name"
                      value={profileData.first_name}
                      onChange={handleProfileChange}
                      className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors"
                      placeholder={t('profile.firstNamePlaceholder')}
                      required
                    />
                  </div>

                  <div>
                    <label htmlFor="last_name" className="block text-sm font-medium text-gray-700 mb-2">
                      {t('users.lastName')}
                    </label>
                    <input
                      type="text"
                      id="last_name"
                      name="last_name"
                      value={profileData.last_name}
                      onChange={handleProfileChange}
                      className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors"
                      placeholder={t('profile.lastNamePlaceholder')}
                      required
                    />
                  </div>

                  <div>
                    <label htmlFor="phone" className="block text-sm font-medium text-gray-700 mb-2">
                      {t('profile.phone')}
                    </label>
                    <input
                      type="tel"
                      id="phone"
                      name="phone"
                      value={profileData.phone}
                      onChange={handleProfileChange}
                      className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors"
                      placeholder={t('profile.phonePlaceholder')}
                    />
                  </div>

                  <div>
                    <label htmlFor="location" className="block text-sm font-medium text-gray-700 mb-2">
                      {t('profile.location')}
                    </label>
                    <input
                      type="text"
                      id="location"
                      name="location"
                      value={profileData.location}
                      onChange={handleProfileChange}
                      className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors"
                      placeholder={t('profile.locationPlaceholder')}
                    />
                  </div>

                  <div>
                    <label htmlFor="website" className="block text-sm font-medium text-gray-700 mb-2">
                      {t('profile.website')}
                    </label>
                    <input
                      type="url"
                      id="website"
                      name="website"
                      value={profileData.website}
                      onChange={handleProfileChange}
                      className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors"
                      placeholder={t('profile.websitePlaceholder')}
                    />
                  </div>

                  <div>
                    <label htmlFor="language" className="block text-sm font-medium text-gray-700 mb-2">
                      {t('profile.language')}
                    </label>
                    <select
                      id="language"
                      name="language"
                      value={profileData.language}
                      onChange={handleProfileChange}
                      className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors"
                    >
                      <option value="en">{t('profile.languageOptions.en')}</option>
                      <option value="es">{t('profile.languageOptions.es')}</option>
                      <option value="fr">{t('profile.languageOptions.fr')}</option>
                      <option value="de">{t('profile.languageOptions.de')}</option>
                      <option value="it">{t('profile.languageOptions.it')}</option>
                      <option value="pt">{t('profile.languageOptions.pt')}</option>
                    </select>
                  </div>
                </div>

                <div>
                  <label htmlFor="bio" className="block text-sm font-medium text-gray-700 mb-2">
                    {t('profile.bio')}
                  </label>
                  <textarea
                    id="bio"
                    name="bio"
                    value={profileData.bio}
                    onChange={handleProfileChange}
                    rows={4}
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors"
                    placeholder={t('profile.bioPlaceholder')}
                  />
                </div>

                <div className="flex justify-end">
                  <button
                    type="submit"
                    disabled={loading}
                    className="bg-primary-600 text-white py-3 px-6 rounded-lg font-medium hover:bg-primary-700 focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    {loading ? (
                      <div className="flex items-center">
                        <Icon name="loading" className="mr-2 animate-spin" />
                        {t('profile.updating')}
                      </div>
                    ) : (
                      t('profile.updateProfile')
                    )}
                  </button>
                </div>
              </form>
            )}

            {/* Change Password Tab */}
            {activeTab === 'password' && (
              <form onSubmit={handlePasswordSubmit} className="space-y-6 max-w-md">
                <div>
                  <label htmlFor="current_password" className="block text-sm font-medium text-gray-700 mb-2">
                    {t('profile.currentPassword')}
                  </label>
                  <input
                    type="password"
                    id="current_password"
                    name="current_password"
                    value={passwordData.current_password}
                    onChange={handlePasswordChange}
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors"
                    placeholder={t('profile.currentPasswordPlaceholder')}
                    required
                  />
                </div>

                <div>
                  <label htmlFor="new_password" className="block text-sm font-medium text-gray-700 mb-2">
                    {t('profile.newPassword')}
                  </label>
                  <input
                    type="password"
                    id="new_password"
                    name="new_password"
                    value={passwordData.new_password}
                    onChange={handlePasswordChange}
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors"
                    placeholder={t('profile.newPasswordPlaceholder')}
                    required
                    minLength={8}
                  />
                  <p className="text-sm text-gray-500 mt-1">{t('profile.errors.passwordMinLength')}</p>
                </div>

                <div>
                  <label htmlFor="confirm_password" className="block text-sm font-medium text-gray-700 mb-2">
                    {t('profile.confirmNewPassword')}
                  </label>
                  <input
                    type="password"
                    id="confirm_password"
                    name="confirm_password"
                    value={passwordData.confirm_password}
                    onChange={handlePasswordChange}
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors"
                    placeholder={t('profile.confirmPasswordPlaceholder')}
                    required
                    minLength={8}
                  />
                </div>

                <div className="flex justify-end">
                  <button
                    type="submit"
                    disabled={loading}
                    className="bg-primary-600 text-white py-3 px-6 rounded-lg font-medium hover:bg-primary-700 focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    {loading ? (
                      <div className="flex items-center">
                        <Icon name="loading" className="mr-2 animate-spin" />
                        {t('profile.changing')}
                      </div>
                    ) : (
                      t('profile.changePassword')
                    )}
                  </button>
                </div>
              </form>
            )}
    </TabLayout>
  );
};

export default Profile;