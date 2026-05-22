import { describe, it, expect, vi, beforeEach } from 'vitest';

const mockOnOk = vi.fn();
const mockOnCancel = vi.fn();
const mockOnDismiss = vi.fn();

vi.mock('quasar', () => ({
  Dialog: {
    create: vi.fn(() => ({
      onOk: (fn: () => void) => { mockOnOk.mockImplementation(() => { fn(); }); return { onCancel: (fn2: () => void) => { mockOnCancel.mockImplementation(() => { fn2(); }); return { onDismiss: (fn3: () => void) => { mockOnDismiss.mockImplementation(() => { fn3(); }); } }; } }; },
    })),
  },
}));

import { dialogConfirm, dialogDelete, dialogWarning, dialogInfo, dialogError } from './dialog';
import { Dialog } from 'quasar';

describe('dialog utils', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('dialogConfirm calls Dialog.create', () => {
    void dialogConfirm({ title: 'Test', message: 'Sure?' });
    expect(Dialog.create).toHaveBeenCalled();
  });

  it('dialogDelete calls Dialog.create with negative ok', () => {
    void dialogDelete({ title: 'Delete', message: 'Delete item?' });
    expect(Dialog.create).toHaveBeenCalledWith(expect.objectContaining({
      title: 'Delete',
      ok: expect.objectContaining({ color: 'negative' }),
    }));
  });

  it('dialogWarning uses warning color', () => {
    void dialogWarning({ title: 'Warn', message: 'Watch out' });
    expect(Dialog.create).toHaveBeenCalledWith(expect.objectContaining({
      ok: expect.objectContaining({ color: 'warning' }),
    }));
  });

  it('dialogInfo uses primary color', () => {
    void dialogInfo({ title: 'Info', message: 'FYI' });
    expect(Dialog.create).toHaveBeenCalledWith(expect.objectContaining({
      ok: expect.objectContaining({ color: 'primary' }),
    }));
  });

  it('dialogError uses negative color', () => {
    void dialogError({ title: 'Error', message: 'Oops' });
    expect(Dialog.create).toHaveBeenCalledWith(expect.objectContaining({
      ok: expect.objectContaining({ color: 'negative' }),
    }));
  });

  it('dialogConfirm resolves true on ok', async () => {
    const promise = dialogConfirm({ title: 'T', message: 'M' });
    mockOnOk();
    await expect(promise).resolves.toBe(true);
  });

  it('dialogConfirm resolves false on cancel', async () => {
    const promise = dialogConfirm({ title: 'T', message: 'M' });
    mockOnCancel();
    await expect(promise).resolves.toBe(false);
  });
});
