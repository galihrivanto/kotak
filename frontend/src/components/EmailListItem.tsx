import React from 'react';
import { EmailListItemProps } from '../types';

const EmailListItem: React.FC<EmailListItemProps> = ({ email, isSelected, onClick }) => {
  const date = new Date(email.received_at).toLocaleString();
  
  return (
    <div 
      className={`email-item ${isSelected ? 'active' : ''}`} 
      onClick={onClick}
    >
      <strong>{email.subject || '(No Subject)'}</strong>
      <div>From: {email.from}</div>
      <div><small>{date}</small></div>
    </div>
  );
};

export default EmailListItem;