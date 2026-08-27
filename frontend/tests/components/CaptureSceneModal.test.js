import { describe, it, expect, vi } from 'vitest';
import { render, fireEvent } from '@testing-library/svelte';
import CaptureSceneModal from '../../src/lib/components/CaptureSceneModal.svelte';

describe('CaptureSceneModal.svelte', () => {
  it('renders scene capture header and icon picker options', () => {
    const { getByText } = render(CaptureSceneModal, {
      props: {
        isOpen: true
      }
    });

    expect(getByText('Capture Multi-Strip Scene')).toBeDefined();
    expect(getByText('Preset Icon')).toBeDefined();
    expect(getByText('Movie')).toBeDefined();
    expect(getByText('Gaming')).toBeDefined();
  });

  it('submits scene capture form with entered name and selected icon', async () => {
    const handleCapture = vi.fn();
    const { getByPlaceholderText, getByText } = render(CaptureSceneModal, {
      props: {
        isOpen: true,
        onCapture: handleCapture
      }
    });

    const input = getByPlaceholderText('e.g. Cyberpunk Ambient');
    await fireEvent.input(input, { target: { value: 'Late Night Chill' } });

    const submitBtn = getByText('Save Scene Snapshot');
    await fireEvent.click(submitBtn);

    expect(handleCapture).toHaveBeenCalledWith('Late Night Chill', 'Sparkles');
  });
});
