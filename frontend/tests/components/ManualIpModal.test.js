import { describe, it, expect, vi } from 'vitest';
import { render, fireEvent } from '@testing-library/svelte';
import ManualIpModal from '../../src/lib/components/ManualIpModal.svelte';

describe('ManualIpModal.svelte', () => {
  it('renders modal header and input', () => {
    const { getByText, getByPlaceholderText } = render(ManualIpModal, {
      props: {
        isOpen: true
      }
    });

    expect(getByText('Add WLED Strip Manually')).toBeDefined();
    expect(getByPlaceholderText('192.168.1.100')).toBeDefined();
  });

  it('submits entered IP address when form is submitted', async () => {
    const handleSubmit = vi.fn();
    const { getByPlaceholderText, getByText } = render(ManualIpModal, {
      props: {
        isOpen: true,
        onAdd: handleSubmit
      }
    });

    const input = getByPlaceholderText('192.168.1.100');
    await fireEvent.input(input, { target: { value: '192.168.1.175' } });

    const submitBtn = getByText('Connect Strip');
    await fireEvent.click(submitBtn);

    expect(handleSubmit).toHaveBeenCalledWith('192.168.1.175');
  });
});
