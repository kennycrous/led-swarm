import { describe, it, expect, vi } from 'vitest';
import { render, fireEvent } from '@testing-library/svelte';
import CaptureSceneModal from '../../src/lib/components/CaptureSceneModal.svelte';

describe('CaptureSceneModal.svelte', () => {
  it('renders scene capture modal when isOpen is true', () => {
    const { getByText } = render(CaptureSceneModal, {
      isOpen: true,
      onClose: () => {},
      onCapture: () => {}
    });

    expect(getByText('Capture Scene Preset')).toBeDefined();
    expect(getByText('Preset Icon')).toBeDefined();
    expect(getByText('Sparkles')).toBeDefined();
  });

  it('submits scene capture form with entered name and selected icon', async () => {
    const onCapture = vi.fn();
    const { getByPlaceholderText, getByText } = render(CaptureSceneModal, {
      isOpen: true,
      onClose: () => {},
      onCapture
    });

    const input = getByPlaceholderText('e.g. Cyberpunk Neon, Movie Time, Relax');
    await fireEvent.input(input, { target: { value: 'Late Night Chill' } });

    const submitBtn = getByText('Snapshot Scene');
    await fireEvent.click(submitBtn);

    expect(onCapture).toHaveBeenCalledWith('Late Night Chill', 'Sparkles', 'global', '');
  });
});
