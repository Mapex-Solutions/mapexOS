import { Response } from 'express';
/**
 * Response to client
 */
export interface ResponseToClient {
    res: Response;
    code: number;
    message: string | string[];
    data: any;
}
