import React from 'react';
import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { MemoryRouter } from 'react-router-dom';
import HomePage from './HomePage';
import apiClient from '../api/apiClient';
import * as tmdbApi from '../api/tmdb';

// Mock API clients
vi.mock('../api/apiClient', () => ({
    default: {
        get: vi.fn(),
    }
}));

vi.mock('../api/tmdb', () => ({
    getMovieVideos: vi.fn(),
}));

describe('HomePage', () => {
    const mockMovies = [
        { id: 1, title: 'Hero Movie 1', overview: 'Overview 1', backdrop_path: '/b1.jpg', poster_path: '/p1.jpg', release_date: '2023-01-01', vote_average: 8.5 },
        { id: 2, title: 'Trending Movie 2', overview: 'Overview 2', backdrop_path: '/b2.jpg', poster_path: '/p2.jpg', release_date: '2023-02-01', vote_average: 7.2 },
        { id: 3, title: 'Trending Movie 3', overview: 'Overview 3', backdrop_path: '/b3.jpg', poster_path: '/p3.jpg', release_date: '2023-03-01', vote_average: 6.9 },
    ];

    beforeEach(() => {
        vi.clearAllMocks();
        vi.useFakeTimers({ shouldAdvanceTime: true });
    });

    afterEach(() => {
        vi.useRealTimers();
    });

    const renderHomePage = () => {
        return render(
            <MemoryRouter>
                <HomePage />
            </MemoryRouter>
        );
    };

    it('renders loading state initially', () => {
        apiClient.get.mockImplementation(() => new Promise(() => { })); // Never resolves
        renderHomePage();
        expect(document.querySelector('.loader')).toBeInTheDocument();
    });

    it('loads and displays popular movies', async () => {
        apiClient.get.mockResolvedValueOnce({ data: { success: true, data: { results: mockMovies } } });
        tmdbApi.getMovieVideos.mockResolvedValueOnce({ data: { success: true, results: [] } });

        renderHomePage();

        await waitFor(() => {
            expect(screen.getByText('Popular on FrameRate')).toBeInTheDocument();
        });

        // Check Hero movie
        expect(screen.getByText('Hero Movie 1')).toBeInTheDocument();
        expect(screen.getByText('Overview 1')).toBeInTheDocument();

        // Check Trending movie
        expect(screen.getByText('Trending Movie 2')).toBeInTheDocument();
    });

    it('fetches and sets the trailer key for the hero movie', async () => {
        apiClient.get.mockResolvedValueOnce({ data: { success: true, data: { results: mockMovies } } });
        tmdbApi.getMovieVideos.mockResolvedValueOnce({
            success: true,
            data: {
                results: [{ type: 'Trailer', site: 'YouTube', key: 'test_trailer_key' }]
            }
        });

        // Mock window.open
        const windowOpenSpy = vi.spyOn(window, 'open').mockImplementation(() => { });

        renderHomePage();

        await waitFor(() => {
            expect(screen.getByText('Hero Movie 1')).toBeInTheDocument();
        });

        // Wait for the button to become enabled
        const watchTrailerBtn = await screen.findByText(/Watch Trailer/i);
        await waitFor(() => {
            expect(watchTrailerBtn).not.toBeDisabled();
        });

        fireEvent.click(watchTrailerBtn);

        expect(windowOpenSpy).toHaveBeenCalledWith('https://www.youtube.com/watch?v=test_trailer_key', '_blank');
        windowOpenSpy.mockRestore();
    });

    it('rotates hero movie on interval', async () => {
        apiClient.get.mockResolvedValue({ data: { success: true, data: { results: mockMovies } } });
        tmdbApi.getMovieVideos.mockResolvedValue({ success: true, data: { results: [] } });

        renderHomePage();

        await waitFor(() => {
            expect(screen.getByText('Hero Movie 1')).toBeInTheDocument();
        });

        // Fast-forward 8 seconds
        vi.advanceTimersByTime(8000);

        // Hero should now be Movie 2
        await waitFor(() => {
            expect(screen.getByText('Trending Movie 2')).toBeInTheDocument();
        });
    });
});
