import React from 'react';
import { AccountSectionProps } from '../types';
import { Icon } from '@iconify/react';

const AccountSection: React.FC<AccountSectionProps> = ({ 
  currentAccount, 
  generateEmail, 
  isLoading 
}) => {
  const copyToClipboard = (text: string): void => {
    navigator.clipboard.writeText(text)
      .then(() => alert('Email address copied to clipboard!'))
      .catch(err => console.error('Failed to copy: ', err));
  };

  return (
    <div className="card">
      <div id="account-section">
        <h2>Your Temporary Email</h2>
        <div id="email-display">
          {currentAccount ? (
            <>
            <div className="email-container">
              <p className="email-address">{currentAccount.email}</p>
              <button 
                className="copy-btn" 
                onClick={() => copyToClipboard(currentAccount.email)}
                disabled={isLoading}
              >
                <Icon icon="mdi:content-copy" style={{ fontSize: '1.25rem' }} />
              </button>
            </div>
            <p><small>This email will be active for this session only.</small></p>
          </>
          ) : (
            <>
              <p>You don't have a temporary email yet. Generate one to get started.</p>
              <button 
                onClick={generateEmail} 
                disabled={isLoading}
              >
                {isLoading ? 'Generating...' : 'Generate Email'}
              </button>
            </>
          )}
        </div>
      </div>
    </div>
  );
};

export default AccountSection;