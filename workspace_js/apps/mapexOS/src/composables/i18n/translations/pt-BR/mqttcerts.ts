export const mqttcerts = {
	tab: { title: 'Certificados' },
	current: {
		title: 'Certificado ativo',
		empty: 'Sem certificado ativo. Clique em Gerar para emitir um novo.',
		serial: 'Serial',
		fingerprint: 'Fingerprint',
		subjectCN: 'CN do sujeito',
		issuedAt: 'Emitido em',
		expiresAt: 'Expira em',
	},
	actions: {
		generate: 'Gerar certificado',
		revoke: 'Revogar',
	},
	dialog: {
		generate: {
			title: 'Gerar certificado MQTT',
			warning:
				'A chave privada é exibida apenas uma vez. Baixe o arquivo e armazene em local seguro; a plataforma nunca a guarda.',
			confirmButton: 'Gerar e baixar',
		},
		replace: {
			title: 'Substituir certificado existente?',
			body: 'Este ativo já possui um certificado ativo. Gerar um novo revoga o atual imediatamente.',
			confirmButton: 'Substituir',
		},
		revoke: {
			title: 'Revogar certificado?',
			body: 'O ativo perderá acesso ao MQTT até que um novo certificado seja emitido.',
			confirmButton: 'Revogar',
		},
	},
	revoked: {
		title: 'Certificados revogados',
		retentionNotice:
			'Certificados revogados são retidos por 30 dias para auditoria. Depois disso, os dados migram para o arquivo de longo prazo (futuro).',
		empty: 'Sem certificados revogados para este ativo.',
		columns: {
			serial: 'Serial',
			reason: 'Motivo',
			revokedAt: 'Revogado em',
		},
	},
	errors: {
		caNotReady: 'Subsistema PKI indisponível. Tente novamente em instantes.',
		replaceRequired: 'Confirmação necessária para substituir o certificado existente.',
		generic: 'Operação falhou. Verifique os logs para detalhes.',
	},
	success: {
		issued: 'Certificado emitido. Download iniciado.',
		revoked: 'Certificado revogado.',
	},
};
