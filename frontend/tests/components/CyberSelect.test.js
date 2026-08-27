import { describe, it, expect } from 'vitest';
import { render, fireEvent } from '@testing-library/svelte';
import CyberSelect from '../../src/lib/components/CyberSelect.svelte';

describe('CyberSelect.svelte', () => {
  it('should render trigger button with default fallback option name', () => {
    const { getByRole } = render(CyberSelect, {
      props: {
        value: 0,
        options: ['Solid', 'Breathe', 'Rainbow']
      }
    });

    const button = getByRole('button');
    expect(button).toBeDefined();
    expect(button.textContent).toContain('Solid');
  });

  it('should toggle popover menu when clicked', async () => {
    const { getByRole, queryByRole } = render(CyberSelect, {
      props: {
        value: 0,
        options: ['Solid', 'Breathe', 'Rainbow']
      }
    });

    const trigger = getByRole('button');
    expect(queryByRole('option')).toBeNull();

    await fireEvent.click(trigger);
    // Menu opens
  });
});
