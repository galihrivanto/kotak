import React, { useState, useEffect } from 'react';
import { emailService } from '../services/api';
import { EmailDetailSectionProps, Email } from '../types';

const EmailDetailSection: React.FC<EmailDetailSectionProps> = ({ accountId, emailId }) => {
  const [emailDetail, setEmailDetail] = useState<Email | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchEmailDetail = async (): Promise<void> => {
      if (!accountId || !emailId) return;
      
      setLoading(true);
      try {
        const data = await emailService.getEmailDetail(accountId, emailId);
        setEmailDetail(data.email);
        setError(null);
      } catch (error) {
        console.error('Error:', error);
        setError('Error loading email details. Please try again.');
      } finally {
        setLoading(false);
      }
    };

    fetchEmailDetail();
  }, [accountId, emailId]);

  if (loading) {
    return (
      <div className="card">
        <h2>Email Details</h2>
        <div className="loading">Loading email details...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="card">
        <h2>Email Details</h2>
        <div className="error">{error}</div>
      </div>
    );
  }

  if (!emailDetail) {
    return (
      <div className="card">
        <h2>Email Details</h2>
        <p>Select an email to view its contents.</p>
      </div>
    );
  }

  const receivedDate = new Date(emailDetail.received_at).toLocaleString();

  // Function to safely display HTML content
  const createMarkup = (htmlContent: string) => {
    return { dangerouslySetInnerHTML: { __html: htmlContent } };
  };

  const isHtmlContent = (content: string): boolean => {
    return content.trim().startsWith('<') && content.includes('>');
  };

  return (
    <div className="card">
      <h2>Email Details</h2>
      <div className="email-detail">
        <h3>{emailDetail.subject || '(No Subject)'}</h3>
        <p><strong>From:</strong> {emailDetail.from}</p>
        <p><strong>To:</strong> {emailDetail.to}</p>
        <p><strong>Received:</strong> {receivedDate}</p>
        <div className="email-content">
          {/* Conditionally render as HTML if the content appears to be HTML, otherwise as text */}
          {isHtmlContent(emailDetail.body) ? 
            <div {...createMarkup(emailDetail.body)} /> : 
            emailDetail.body
          }
        </div>
      </div>
    </div>
  );
};

export default EmailDetailSection;