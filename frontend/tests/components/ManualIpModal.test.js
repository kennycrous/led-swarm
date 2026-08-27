import { describe, it, expect, vi } from 'vitest';
import { render, fireEvent } from '@testing-library/svelte';
import ManualIpModal from '../../src/lib/components/ManualIpModal.svelte';

describe('ManualIpModal.svelte', () => {
  it('renders modal when isOpen is true', () => {
    const { getByText } = render(ManualIpModal, {
      isOpen: true,
      onClose: () => {},
      onAdd: () => {}
    });

    expect(getByText('Add WLED Device by IP')).toBeDefined();
  });

  it('submits typed IP address', async () => {
    const onAdd = vi.fn().mockResolvedValue({});
    const { getByPlaceholderText, getByText } = render(ManualIpModal, {
      isOpen: true,
      onClose: () => {},
      onAdd
    });

    const input = getByPlaceholderText('192.168.1.150');
    await fireEvent.input(input, { target: { value: '192.168.1.175' } });

    const submitBtn = getByText('ADD DEVICE');
    await fireEvent.click(submitBtn);

    expect(onAdd).toHaveBeenCalledWith('192.168.1.175');
  });
});
