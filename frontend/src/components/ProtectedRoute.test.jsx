import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { MemoryRouter, Routes, Route } from 'react-router-dom';
import ProtectedRoute from './ProtectedRoute';
import * as AuthContext from '../context/AuthContext';

// Mock AuthContext
vi.mock('../context/AuthContext', () => ({
    useAuth: vi.fn(),
}));

describe('ProtectedRoute', () => {
    const TestChild = () => <div data-testid="protected-content">Protected Content</div>;
    const LoginRedirect = () => <div data-testid="login-redirect">Login Page</div>;
    const HomeRedirect = () => <div data-testid="home-redirect">Home Page</div>;

    const renderWithRouter = (ui, initialEntry = '/') => {
        return render(
            <MemoryRouter initialEntries={[initialEntry]}>
                <Routes>
                    <Route path="/" element={<HomeRedirect />} />
                    <Route path="/login" element={<LoginRedirect />} />
                    <Route path="/protected" element={ui} />
                </Routes>
            </MemoryRouter>
        );
    };

    it('shows loading spinner when auth is loading', () => {
        vi.spyOn(AuthContext, 'useAuth').mockReturnValue({ loading: true, user: null });

        renderWithRouter(
            <ProtectedRoute>
                <TestChild />
            </ProtectedRoute>,
            '/protected'
        );

        expect(screen.queryByTestId('protected-content')).not.toBeInTheDocument();
        // Use container query for the loader class
        expect(document.querySelector('.loader')).toBeInTheDocument();
    });

    it('redirects to login when user is not authenticated', () => {
        vi.spyOn(AuthContext, 'useAuth').mockReturnValue({ loading: false, user: null });

        renderWithRouter(
            <ProtectedRoute>
                <TestChild />
            </ProtectedRoute>,
            '/protected'
        );

        expect(screen.queryByTestId('protected-content')).not.toBeInTheDocument();
        expect(screen.getByTestId('login-redirect')).toBeInTheDocument();
    });

    it('renders children when user is authenticated', () => {
        vi.spyOn(AuthContext, 'useAuth').mockReturnValue({
            loading: false,
            user: { id: 1, username: 'user', is_admin: false }
        });

        renderWithRouter(
            <ProtectedRoute>
                <TestChild />
            </ProtectedRoute>,
            '/protected'
        );

        expect(screen.getByTestId('protected-content')).toBeInTheDocument();
        expect(screen.queryByTestId('login-redirect')).not.toBeInTheDocument();
    });

    it('redirects to home when admin is required but user is not admin', () => {
        vi.spyOn(AuthContext, 'useAuth').mockReturnValue({
            loading: false,
            user: { id: 1, username: 'user', is_admin: false }
        });

        renderWithRouter(
            <ProtectedRoute requireAdmin={true}>
                <TestChild />
            </ProtectedRoute>,
            '/protected'
        );

        expect(screen.queryByTestId('protected-content')).not.toBeInTheDocument();
        expect(screen.getByTestId('home-redirect')).toBeInTheDocument();
    });

    it('renders children when admin is required and user is admin', () => {
        vi.spyOn(AuthContext, 'useAuth').mockReturnValue({
            loading: false,
            user: { id: 1, username: 'admin', is_admin: true }
        });

        renderWithRouter(
            <ProtectedRoute requireAdmin={true}>
                <TestChild />
            </ProtectedRoute>,
            '/protected'
        );

        expect(screen.getByTestId('protected-content')).toBeInTheDocument();
        expect(screen.queryByTestId('home-redirect')).not.toBeInTheDocument();
    });
});
