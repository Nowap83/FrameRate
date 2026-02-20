import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';

describe('Basic Test Setup', () => {
    it('renders a react element', () => {
        render(<div>Hello Vitest</div>);
        expect(screen.getByText('Hello Vitest')).toBeInTheDocument();
    });
});
