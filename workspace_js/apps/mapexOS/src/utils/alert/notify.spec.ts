import { describe, it, expect, vi, beforeEach } from 'vitest';

// Mock Quasar Notify before importing
vi.mock('quasar', () => ({
  Notify: {
    create: vi.fn(() => vi.fn()),
  },
}));

import { notifySuccess, notifyInfo, notifyFail, notifyWarning } from './notify';
import { Notify } from 'quasar';

describe('notify utils', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.useFakeTimers();
  });

  it('notifySuccess calls Notify.create with positive color', () => {
    notifySuccess({ message: 'Done!' });
    vi.advanceTimersByTime(400);
    expect(Notify.create).toHaveBeenCalledWith(expect.objectContaining({
      color: 'positive',
      icon: 'check',
      message: 'Done!',
    }));
  });

  it('notifyInfo calls Notify.create with primary color', () => {
    notifyInfo({ message: 'Info' });
    vi.advanceTimersByTime(400);
    expect(Notify.create).toHaveBeenCalledWith(expect.objectContaining({
      color: 'primary',
      icon: 'info',
    }));
  });

  it('notifyFail calls Notify.create with red color', () => {
    notifyFail({ message: 'Error' });
    vi.advanceTimersByTime(400);
    expect(Notify.create).toHaveBeenCalledWith(expect.objectContaining({
      color: 'red-4',
      icon: 'report_problem',
    }));
  });

  it('notifyWarning calls Notify.create with warning color', () => {
    notifyWarning({ message: 'Warning' });
    vi.advanceTimersByTime(400);
    expect(Notify.create).toHaveBeenCalledWith(expect.objectContaining({
      color: 'warning',
      icon: 'warning',
    }));
  });

  it('passes html option when provided', () => {
    notifySuccess({ message: '<b>Bold</b>', html: true });
    vi.advanceTimersByTime(400);
    expect(Notify.create).toHaveBeenCalledWith(expect.objectContaining({
      html: true,
    }));
  });

  it('passes timeout option when provided', () => {
    notifySuccess({ message: 'Quick', timeout: 1000 });
    vi.advanceTimersByTime(400);
    expect(Notify.create).toHaveBeenCalledWith(expect.objectContaining({
      timeout: 1000,
    }));
  });
});
