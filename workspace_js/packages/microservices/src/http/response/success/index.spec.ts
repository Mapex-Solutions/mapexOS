import { Response } from 'express';
import { success, created } from './';

// Global vars
let responseObject: any = {};

/**
 * Mock response HTTP express
 */
const mockResponse = {} as unknown as Response;
mockResponse.status = jest.fn().mockImplementation(() => mockResponse),
  mockResponse.json = jest.fn().mockImplementation(result => {
    responseObject = result;
  });

/**
 * Start test functionality
 */
describe('Unit test rejection', () => {

  it('Success - 200', async () => {
    // Arrange
    const data = { name: 'Mapex', age: 23 };

    // Act
    await success(mockResponse, data);

    // Assert
    expect(responseObject).toEqual(
      expect.objectContaining({
        httpStatus: 200,
        error: null,
        data: { name: 'Mapex', age: 23 }
      })
    );
  });

  it('Created - 201', async () => {
    // Arrange
    const data = { name: 'Mapex', age: 23 };

    // Act
    await created(mockResponse, data);

    // Assert
    expect(responseObject).toEqual(
      expect.objectContaining({
        httpStatus: 201,
        error: null,
        data: { name: 'Mapex', age: 23 }
      })
    );
  });
});