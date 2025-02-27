// Account information
export interface Account {
    account_id: string;
    email: string;
    created_at: string;
  }
  
  // Email message structure
  export interface Email {
    id: string;
    from: string;
    to: string;
    subject: string;
    body: string;
    received_at: string;
  }
  
  // API response for emails list
  export interface EmailsResponse {
    emails: Email[];
  }
  
  // API response for single email
  export interface EmailDetailResponse {
    email: Email;
  }
  
  // Component props
  export interface AccountSectionProps {
    currentAccount: Account | null;
    generateEmail: () => Promise<void>;
    isLoading: boolean;
  }
  
  export interface InboxSectionProps {
    emails: Email[];
    onRefresh: () => Promise<void>;
    onSelectEmail: (emailId: string) => void;
    selectedEmailId: string | null;
    isLoading: boolean;
  }
  
  export interface EmailListItemProps {
    email: Email;
    isSelected: boolean;
    onClick: () => void;
  }
  
  export interface EmailDetailSectionProps {
    accountId: string | undefined;
    emailId: string | null;
  }