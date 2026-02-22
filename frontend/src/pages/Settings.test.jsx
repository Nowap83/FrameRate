import React from 'react';
import { render, screen, waitFor, act, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { MemoryRouter } from 'react-router-dom';
import Settings from './Settings';
import apiClient from '../api/apiClient';
import * as AuthContext from '../context/AuthContext';

// Mock apiClient
vi.mock('../api/apiClient', () => ({
    default: {
        get: vi.fn(),
        put: vi.fn(),
        post: vi.fn(),
    }
}));

// Mock AuthContext
vi.mock('../context/AuthContext', () => ({
    useAuth: vi.fn(),
}));

// Mock MovieSearch component
vi.mock('../components/MovieSearch', () => ({
    default: ({ onSelect, onCancel }) => (
        <div data-testid="mock-movie-search">
            <button onClick={() => onSelect({ id: 999, title: 'Mock Movie', poster_path: '/mock.jpg', release_date: '2023-01-01' })}>Select Mock Movie</button>
            <button onClick={onCancel}>Cancel Search</button>
        </div>
    )
}));

describe('Settings Page', () => {
    const mockUser = {
        id: 1,
        username: 'testuser',
        email: 'test@example.com',
        profile_picture_url: '/avatar.jpg'
    };

    const mockProfileData = {
        user: {
            username: 'testuser',
            bio: 'test bio',
            given_name: 'Test',
            family_name: 'User',
            location: 'Earth',
            website: 'https://test.com',
            profile_picture_url: '/avatar.jpg'
        },
        favorites: [
            { tmdb_id: 100, title: 'Fav 1', poster_url: '/fav1.jpg' }
        ]
    };

    beforeEach(() => {
        vi.clearAllMocks();
    });

    const renderSettings = () => {
        return render(
            <MemoryRouter>
                <Settings />
            </MemoryRouter>
        );
    };

    it('renders loading state initially', () => {
        vi.spyOn(AuthContext, 'useAuth').mockReturnValue({ user: null, loading: true, setUser: vi.fn() });
        renderSettings();
        expect(document.querySelector('.loader')).toBeInTheDocument();
    });

    it('loads profile data and renders correctly', async () => {
        vi.spyOn(AuthContext, 'useAuth').mockReturnValue({ user: mockUser, loading: false, setUser: vi.fn() });
        apiClient.get.mockResolvedValueOnce({ data: mockProfileData });

        renderSettings();

        // Verify loading eventually goes away
        await waitFor(() => {
            expect(document.querySelector('.loader')).not.toBeInTheDocument();
        });

        expect(screen.getByDisplayValue('testuser')).toBeInTheDocument();
        expect(screen.getByDisplayValue('test bio')).toBeInTheDocument();
        expect(screen.getByDisplayValue('Test')).toBeInTheDocument();
        expect(screen.getByDisplayValue('Earth')).toBeInTheDocument();

        // Check favorite film loaded
        const images = screen.getAllByRole('img');
        expect(images.some(img => img.alt === 'Fav 1')).toBe(true);
    });

    it('navigates between tabs', async () => {
        vi.spyOn(AuthContext, 'useAuth').mockReturnValue({ user: mockUser, loading: false, setUser: vi.fn() });
        apiClient.get.mockResolvedValueOnce({ data: mockProfileData });

        renderSettings();

        await waitFor(() => {
            expect(screen.getByRole('button', { name: /profile/i })).toBeInTheDocument();
        });

        // Go to Auth tab
        fireEvent.click(screen.getByRole('button', { name: /auth/i }));
        expect(screen.getByText('Change Password')).toBeInTheDocument();

        // Go to Avatar tab
        fireEvent.click(screen.getByRole('button', { name: /avatar/i }));
        expect(screen.getByText('Change Avatar')).toBeInTheDocument();
    });

    it('submits profile updates correctly', async () => {
        vi.spyOn(AuthContext, 'useAuth').mockReturnValue({ user: mockUser, loading: false, setUser: vi.fn() });
        apiClient.get.mockResolvedValueOnce({ data: mockProfileData });
        apiClient.put.mockResolvedValueOnce({ data: { message: 'Success' } });

        renderSettings();

        await waitFor(() => {
            expect(screen.getByDisplayValue('testuser')).toBeInTheDocument();
        });

        // Change given name
        const givenNameInput = screen.getByDisplayValue('Test');
        fireEvent.change(givenNameInput, { target: { name: 'given_name', value: 'Updated Test' } });

        fireEvent.click(screen.getByText('Save Changes'));

        await waitFor(() => {
            expect(apiClient.put).toHaveBeenCalledWith('/users/me', expect.objectContaining({
                given_name: 'Updated Test',
                username: 'testuser'
                // further fields could be checked
            }));
        });
    });

    it('checks username availability', async () => {
        vi.spyOn(AuthContext, 'useAuth').mockReturnValue({ user: mockUser, loading: false, setUser: vi.fn() });
        apiClient.get.mockResolvedValueOnce({ data: mockProfileData });
        // Mock check-username response
        apiClient.get.mockResolvedValueOnce({ data: { available: true } });

        renderSettings();

        await waitFor(() => {
            expect(screen.getByDisplayValue('testuser')).toBeInTheDocument();
        });

        const usernameInput = screen.getByDisplayValue('testuser');
        fireEvent.change(usernameInput, { target: { name: 'username', value: 'newusername' } });

        const checkButton = screen.getByText('Check');
        fireEvent.click(checkButton);

        await waitFor(() => {
            expect(apiClient.get).toHaveBeenCalledWith('/users/check-username?username=newusername');
            expect(screen.getByText('Username is available!')).toBeInTheDocument();
        });
    });
});
