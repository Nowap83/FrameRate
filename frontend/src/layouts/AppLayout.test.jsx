import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import AppLayout from './AppLayout';

// Mock Header and Footer to isolate testing AppLayout
vi.mock('../components/Header', () => ({
    default: () => <header data-testid="mock-header">Header Mock</header>,
}));

vi.mock('../components/Footer', () => ({
    default: () => <footer data-testid="mock-footer">Footer Mock</footer>,
}));

describe('AppLayout Component', () => {
    it('renders Header, children, and Footer in correct order', () => {
        render(
            <AppLayout>
                <div data-testid="test-children">Main Content Area</div>
            </AppLayout>
        );

        // Verify Header is rendered
        expect(screen.getByTestId('mock-header')).toBeInTheDocument();

        // Verify Children are rendered inside main content area
        expect(screen.getByTestId('test-children')).toBeInTheDocument();
        expect(screen.getByText('Main Content Area')).toBeInTheDocument();

        // Verify Footer is rendered
        expect(screen.getByTestId('mock-footer')).toBeInTheDocument();
    });
});
