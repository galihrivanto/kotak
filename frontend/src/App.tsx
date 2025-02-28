import React, { useState, useEffect, useCallback } from 'react';
import { emailService } from './services/api';
import AccountSection from './components/AccountSection';
import InboxSection from './components/InboxSection';
import EmailDetailSection from './components/EmailDetailSection';
import { Account, Email } from './types';
import './styles.css';

const App: React.FC = () => {
    const [currentAccount, setCurrentAccount] = useState<Account | null>(null);
    const [currentEmails, setCurrentEmails] = useState<Email[]>([]);
    const [selectedEmailId, setSelectedEmailId] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);

    const checkAccount = useCallback(async (accountId: string): Promise<void> => {
        const exists = await emailService.checkAccount(accountId);
        if (!exists) {
            setCurrentAccount(null);
            localStorage.removeItem('tempEmailAccount');
        }
    }, []);

    // Check for existing account in localStorage on component mount
    useEffect(() => {
        const savedAccount = localStorage.getItem('tempEmailAccount');
        if (savedAccount) {
            try {
                const account = JSON.parse(savedAccount);
                setCurrentAccount(account);
                checkAccount(account.account_id); 
            } catch (e) {
                console.error('Error parsing saved account:', e);
                localStorage.removeItem('tempEmailAccount');
            }
        }
    }, [checkAccount]);

    // Fetch emails for the current account
    const fetchEmails = useCallback(async (): Promise<void> => {
        if (!currentAccount) return;
        setIsLoading(true);
        try {
            const data = await emailService.getEmails(currentAccount.account_id);
            setCurrentEmails(data.emails || []);
            setError(null);
        } catch (error) {
            console.error('Error:', error);
            setError('Failed to fetch emails. Please try refreshing.');
        } finally {
            setIsLoading(false);
        }
    }, [currentAccount]);

    // Fetch emails when account exists
    useEffect(() => {
        if (currentAccount) {
            fetchEmails();
            // Set up auto-refresh
            const intervalId = setInterval(fetchEmails, 30000);
            return () => clearInterval(intervalId);
        }
    }, [currentAccount, fetchEmails]);

    // Generate new email address
    const generateEmail = async (): Promise<void> => {
        setIsLoading(true);
        setError(null);

        try {
            const account = await emailService.generateEmailAccount();
            setCurrentAccount(account);
            localStorage.setItem('tempEmailAccount', JSON.stringify(account));
        } catch (error) {
            console.error('Error:', error);
            setError('Failed to generate email address. Please try again.');
        } finally {
            setIsLoading(false);
        }
    };

    // View a specific email
    const viewEmail = (emailId: string): void => {
        setSelectedEmailId(emailId);
    };

    return (
        <div className="container">
            <div className="main-header">
                <img src="/box.svg" alt="Kotak" width={64} height={64} />
                <h1>Kotak | Temporary Email</h1>
            </div>

            {error && <div className="error-banner">{error}</div>}

            <AccountSection
                currentAccount={currentAccount}
                generateEmail={generateEmail}
                isLoading={isLoading}
            />

            {currentAccount && (
                <InboxSection
                    emails={currentEmails}
                    onRefresh={fetchEmails}
                    onSelectEmail={viewEmail}
                    selectedEmailId={selectedEmailId}
                    isLoading={isLoading}
                />
            )}

            {selectedEmailId && (
                <EmailDetailSection
                    accountId={currentAccount?.account_id}
                    emailId={selectedEmailId}
                />
            )}
        </div>
    );
};

export default App;