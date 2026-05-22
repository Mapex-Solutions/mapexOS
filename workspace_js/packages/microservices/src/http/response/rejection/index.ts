import { HttpStatus } from '@src/http';
import { Response } from 'express';

export const badRequest = async (res: Response, error: string | string[]) => {
	const httpStatus: number = HttpStatus.BAD_REQUEST;

	return res
		.status(httpStatus)
		.json({ httpStatus, error, data: null });
};

export const internalError = async (res: Response, error: string | string[]) => {
	const httpStatus: number = HttpStatus.INTERNAL_SERVER_ERROR;

	return res
		.status(httpStatus)
		.json({ httpStatus, error, data: null });
};

export const unauthorized = async (res: Response, error: string | string[]) => {
	const httpStatus: number = HttpStatus.UNAUTHORIZED;

	return res
		.status(httpStatus)
		.json({ httpStatus, error, data: null });
};

export const forbidden = async (res: Response, error: string | string[]) => {
	const httpStatus: number = HttpStatus.FORBIDDEN;

	return res
		.status(httpStatus)
		.json({ httpStatus, error, data: null });
};

export const notFound = async (res: Response, error: string | string[]) => {
	const httpStatus: number = HttpStatus.NOT_FOUND;

	return res
		.status(httpStatus)
		.json({ httpStatus, error, data: null });
};

export const conflict = async (res: Response, error: string | string[]) => {
	const httpStatus: number = HttpStatus.CONFLICT;

	return res
		.status(httpStatus)
		.json({ httpStatus, error, data: null });
};