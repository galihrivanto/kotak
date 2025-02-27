import React from 'react';
import { AccountSectionProps } from '../types';

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
              <p>Your temporary email address is:</p>
              <p className="email-address">{currentAccount.email}</p>
              <button 
                className="copy-btn" 
                onClick={() => copyToClipboard(currentAccount.email)}
                disabled={isLoading}
              >
                Copy to Clipboard
              </button>
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