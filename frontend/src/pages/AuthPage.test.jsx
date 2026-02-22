import React from 'react';
import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { MemoryRouter } from 'react-router-dom';
import AuthPage from './AuthPage';
import { authService } from '../api/auth';
import * as AuthContext from '../context/AuthContext';

// Mock auth service
vi.mock('../api/auth', () => ({
    authService: {
        login: vi.fn(),
        register: vi.fn(),
    }
}));

// Mock AuthContext
const mockLogin = vi.fn();
vi.mock('../context/AuthContext', () => ({
    useAuth: () => ({
        login: mockLogin,
        user: null,
    }),
}));

// Mock framer-motion to skip animations in tests
vi.mock('framer-motion', async () => {
    const actual = await vi.importActual('framer-motion');
    return {
        ...actual,
        AnimatePresence: ({ children }) => <>{children}</>,
        motion: {
            div: ({ children, ...props }) => <div {...props}>{children}</div>,
            h2: ({ children, ...props }) => <h2 {...props}>{children}</h2>,
            p: ({ children, ...props }) => <p {...props}>{children}</p>,
        }
    };
});

describe('AuthPage - Login Mode', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    const renderAuthPage = (initialEntries = ['/login']) => {
        return render(
            <MemoryRouter initialEntries={initialEntries}>
                <AuthPage />
            </MemoryRouter>
        );
    };

    it('renders login form correctly', () => {
        renderAuthPage();
        expect(screen.getByText('Welcome Back')).toBeInTheDocument();
        expect(screen.getByPlaceholderText('Enter your email...')).toBeInTheDocument();
        expect(screen.getByPlaceholderText('••••••••')).toBeInTheDocument();
    });

    it('handles successful login', async () => {
        authService.login.mockResolvedValueOnce({ token: 'fake-token', user: { username: 'test' } });

        renderAuthPage();

        fireEvent.change(screen.getByPlaceholderText('Enter your email...'), { target: { value: 'test@example.com' } });
        fireEvent.change(screen.getByPlaceholderText('••••••••'), { target: { value: 'password123' } });

        fireEvent.submit(screen.getByRole('button', { name: /Login/i }));

        await waitFor(() => {
            expect(authService.login).toHaveBeenCalledWith({ login: 'test@example.com', password: 'password123' });
            expect(mockLogin).toHaveBeenCalledWith({ username: 'test' }, 'fake-token');
        });
    });

    it('handles login failure and shows error', async () => {
        authService.login.mockRejectedValueOnce({
            response: { data: { error: 'Invalid credentials' } }
        });

        renderAuthPage();

        fireEvent.change(screen.getByPlaceholderText('Enter your email...'), { target: { value: 'test@example.com' } });
        fireEvent.change(screen.getByPlaceholderText('••••••••'), { target: { value: 'wrongpassword' } });

        fireEvent.submit(screen.getByRole('button', { name: /Login/i }));

        await waitFor(() => {
            expect(screen.getByText('Invalid credentials')).toBeInTheDocument();
        });
    });
});

describe('AuthPage - Register Mode', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    const renderAuthPage = () => {
        return render(
            <MemoryRouter initialEntries={['/register']}>
                <AuthPage />
            </MemoryRouter>
        );
    };

    it('renders register form correctly', () => {
        renderAuthPage();
        expect(screen.getByText('Create Account')).toBeInTheDocument();
        expect(screen.getByPlaceholderText('name@example.com')).toBeInTheDocument();
    });

    it('handles successful registration', async () => {
        authService.register.mockResolvedValueOnce({});

        renderAuthPage();

        fireEvent.change(screen.getByPlaceholderText('name@example.com'), { target: { value: 'new@example.com' } });
        fireEvent.change(screen.getByPlaceholderText('Choose a username'), { target: { value: 'newuser' } });

        // There are two password fields in register form
        const passwordInputs = screen.getAllByPlaceholderText('••••••••');
        fireEvent.change(passwordInputs[0], { target: { value: 'Password123!' } });
        fireEvent.change(passwordInputs[1], { target: { value: 'Password123!' } });

        fireEvent.submit(screen.getByRole('button', { name: /Register/i }));

        await waitFor(() => {
            expect(authService.register).toHaveBeenCalledWith(expect.objectContaining({
                email: 'new@example.com',
                username: 'newuser',
                password: 'Password123!'
            }));
            expect(screen.getByText('Check your email')).toBeInTheDocument();
        });
    });
});
