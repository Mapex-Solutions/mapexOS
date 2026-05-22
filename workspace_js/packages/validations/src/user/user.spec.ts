// zodSchemas.test.ts
import { Email, Password, IdCard } from './user.validation'; // Ajuste o caminho conforme sua estrutura

describe('Email Schema', () => {
  it('should pass with a valid email', () => {
    const validEmail = 'test@example.com';
    expect(Email.parse(validEmail)).toBe(validEmail);
  });

  it('should fail with an invalid email', () => {
    const invalidEmail = 'invalid-email';
    expect(() => Email.parse(invalidEmail)).toThrow();
  });
});

describe('Password Schema', () => {
  it('should pass with a valid password', () => {
    // A senha válida deve ter ao menos 8 caracteres, conter ao menos uma letra, um número e um caractere especial.
    const validPassword = 'Abc123@!';
    expect(Password.parse(validPassword)).toBe(validPassword);
  });

  it('should fail if the password is less than 8 characters', () => {
    const shortPassword = 'Ab1@';
    expect(() => Password.parse(shortPassword)).toThrow(/Password must have at least 8 characters/);
  });

  it('should fail if the password does not contain any letter', () => {
    const noLetterPassword = '12345678@';
    expect(() => Password.parse(noLetterPassword)).toThrow(/Password must contain at least one letter/);
  });

  it('should fail if the password does not contain any number', () => {
    const noNumberPassword = 'Abcdefg@';
    expect(() => Password.parse(noNumberPassword)).toThrow(/Password must contain at least one number/);
  });

  it('should fail if the password does not contain any special character', () => {
    const noSpecialCharPassword = 'Abc12345';
    expect(() => Password.parse(noSpecialCharPassword)).toThrow(/Password must contain at least one special character/);
  });
});

describe('IdCard Schema (CPF Validation)', () => {
  /**
   * Observação:
   * O CPF "123.456.789-09" é matematicamente válido segundo as regras de cálculo,
   * mesmo que não seja um CPF real emitido. Portanto, para testar CPFs inválidos,
   * utilize outros exemplos.
   */

  it('should pass with a valid CPF with punctuation', () => {
    const validCpfWithPunctuation = '111.444.777-35';
    expect(IdCard.parse(validCpfWithPunctuation)).toBe(validCpfWithPunctuation);
  });

  it('should pass with a valid CPF without punctuation', () => {
    const validCpfWithoutPunctuation = '11144477735';
    expect(IdCard.parse(validCpfWithoutPunctuation)).toBe(validCpfWithoutPunctuation);
  });

  it('should fail with an invalid CPF (all identical digits)', () => {
    const invalidCpf = '111.111.111-11';
    expect(() => IdCard.parse(invalidCpf)).toThrow();
  });

  it('should fail with an invalid CPF (wrong verifier digit)', () => {
    // Alterando o último dígito de um CPF válido para forçar erro na verificação
    const invalidCpf = '111.444.777-34'; // O último dígito deveria ser 5, conforme o CPF válido '111.444.777-35'
    expect(() => IdCard.parse(invalidCpf)).toThrow();
  });
});
