import { Account, EmailsResponse, EmailDetailResponse } from '../types';

const API_BASE_URL = import.meta.env.VITE_API_HOST + import.meta.env.VITE_API_BASE;

export const emailService = {
  // Check if an account exists
  checkAccount: async (accountId: string): Promise<boolean> => {
    const response = await fetch(`${API_BASE_URL}/accounts/${accountId}`);
    return response.ok;
  },

  // Generate a new temporary email account
  generateEmailAccount: async (): Promise<Account> => {
    const response = await fetch(`${API_BASE_URL}/accounts`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      }
    });
    
    if (!response.ok) {
      throw new Error('Failed to generate email account');
    }
    
    return response.json();
  },
  
  // Get all emails for an account
  getEmails: async (accountId: string): Promise<EmailsResponse> => {
    const response = await fetch(`${API_BASE_URL}/accounts/${accountId}/emails`);
    
    if (!response.ok) {
      throw new Error('Failed to fetch emails');
    }
    
    return response.json();
  },
  
  // Get a specific email's details
  getEmailDetail: async (accountId: string, emailId: string): Promise<EmailDetailResponse> => {
    const response = await fetch(`${API_BASE_URL}/accounts/${accountId}/emails/${emailId}`);
    
    if (!response.ok) {
      throw new Error('Failed to fetch email details');
    }
    
    return response.json();
  }
};