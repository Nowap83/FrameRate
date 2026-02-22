import { render, screen, waitFor, act } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { AuthProvider, useAuth } from './AuthContext';
import apiClient from '../api/apiClient';
import React from 'react';

// Mock the apiClient
vi.mock('../api/apiClient', () => ({
    default: {
        get: vi.fn(),
    }
}));

// Test component to consume the context
const TestComponent = () => {
    const { user, login, logout, loading } = useAuth();

    if (loading) return <div>Loading...</div>;

    return (
        <div>
            <div data-testid="user-status">{user ? `Logged in as ${user.username}` : 'Not logged in'}</div>
            <button onClick={() => login({ username: 'testuser' }, 'fake-token')}>Login</button>
            <button onClick={logout}>Logout</button>
        </div>
    );
};

describe('AuthContext', () => {
    beforeEach(() => {
        localStorage.clear();
        vi.clearAllMocks();
    });

    it('provides initial state without token', async () => {
        render(
            <AuthProvider>
                <TestComponent />
            </AuthProvider>
        );

        // After loading, it should not be logged in
        await waitFor(() => {
            expect(screen.getByTestId('user-status')).toHaveTextContent('Not logged in');
        });
    });

    it('loads user data if token is present', async () => {
        localStorage.setItem('token', 'valid-token');

        // Mock successful API response
        apiClient.get.mockResolvedValueOnce({ data: { user: { username: 'apiuser' } } });

        render(
            <AuthProvider>
                <TestComponent />
            </AuthProvider>
        );

        // It should eventually show the logged in user
        await waitFor(() => {
            expect(screen.getByTestId('user-status')).toHaveTextContent('Logged in as apiuser');
        });

        expect(apiClient.get).toHaveBeenCalledWith('/users/me');
    });

    it('clears token if API request fails', async () => {
        localStorage.setItem('token', 'invalid-token');

        // Mock failed API response
        apiClient.get.mockRejectedValueOnce(new Error('Unauthorized'));

        render(
            <AuthProvider>
                <TestComponent />
            </AuthProvider>
        );

        await waitFor(() => {
            expect(screen.getByTestId('user-status')).toHaveTextContent('Not logged in');
        });

        expect(localStorage.getItem('token')).toBeNull();
    });

    it('handles manual login and logout', async () => {
        render(
            <AuthProvider>
                <TestComponent />
            </AuthProvider>
        );

        // Wait for initial load to finish
        await waitFor(() => {
            expect(screen.getByTestId('user-status')).toHaveTextContent('Not logged in');
        });

        // Login
        act(() => {
            screen.getByText('Login').click();
        });

        expect(screen.getByTestId('user-status')).toHaveTextContent('Logged in as testuser');
        expect(localStorage.getItem('token')).toBe('fake-token');

        // Logout
        act(() => {
            screen.getByText('Logout').click();
        });

        expect(screen.getByTestId('user-status')).toHaveTextContent('Not logged in');
        expect(localStorage.getItem('token')).toBeNull();
    });

    it('logs out on auth:unauthorized event', async () => {
        render(
            <AuthProvider>
                <TestComponent />
            </AuthProvider>
        );

        await waitFor(() => {
            expect(screen.getByTestId('user-status')).toHaveTextContent('Not logged in');
        });

        // Manually login first
        act(() => {
            screen.getByText('Login').click();
        });
        expect(screen.getByTestId('user-status')).toHaveTextContent('Logged in as testuser');

        // Dispatch the custom event
        act(() => {
            window.dispatchEvent(new Event('auth:unauthorized'));
        });

        expect(screen.getByTestId('user-status')).toHaveTextContent('Not logged in');
        expect(localStorage.getItem('token')).toBeNull();
    });
});
