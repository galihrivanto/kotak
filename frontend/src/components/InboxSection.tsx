import React from 'react';
import EmailListItem from './EmailListItem';
import { InboxSectionProps } from '../types';

const InboxSection: React.FC<InboxSectionProps> = ({ 
  emails, 
  onRefresh, 
  onSelectEmail, 
  selectedEmailId, 
  isLoading 
}) => {
  return (
    <div className="card">
      <div className="header">
        <h2>Inbox</h2>
        <button 
          onClick={onRefresh} 
          disabled={isLoading}
        >
          {isLoading ? 'Refreshing...' : 'Refresh'}
        </button>
      </div>
      
      <div className="email-list">
        {isLoading && emails.length === 0 ? (
          <div className="loading">Loading emails...</div>
        ) : emails.length === 0 ? (
          <div className="loading">No emails yet. When you receive emails, they will appear here.</div>
        ) : (
          emails.map(email => (
            <EmailListItem 
              key={email.id}
              email={email}
              isSelected={email.id === selectedEmailId}
              onClick={() => onSelectEmail(email.id)}
            />
          ))
        )}
      </div>
    </div>
  );
};

export default InboxSection;