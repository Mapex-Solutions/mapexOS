/**
 * Portuguese (pt-BR) translations for core workflow plugins.
 * Organized by plugin ID → node short name → translation keys.
 *
 * Namespace at runtime: wf.{pluginId}.nodes.{nodeShortName}.*
 */

/** Core Triggers plugin translations */
export const coreTriggers = {
  nodes: {
    start: {
      label: 'Início',
      description: 'Ponto de entrada do workflow',
    },
    trigger_event: {
      label: 'Evento de Gatilho',
      description: 'Ponto de entrada — disparado por um gatilho configurado na plataforma',
      config: {
        triggerSection: 'Gatilho',
        changeTrigger: 'Alterar gatilho',
        removeTrigger: 'Remover gatilho',
        selectTrigger: 'Selecionar Gatilho',
        variablesSection: 'Variáveis',
        selectPrompt: 'Selecione um gatilho para configurar este nó de evento',
      },
    },
  },
};

/** Core Logic plugin translations */
export const coreLogic = {
  nodes: {
    condition: {
      label: 'Condições',
      description: 'Avalia condições e ramifica com base no resultado',
      config: {
        noConditionsYet: 'Nenhuma condição ainda',
        conditionLabel: 'Condição',
        conditionDescription: 'Comparação simples de campo',
        groupLabel: 'Grupo',
        groupDescription: 'Grupo de condições com lógica',
        addButton: 'Adicionar',
        rename: 'Renomear',
        deleteItem: 'Excluir',
        deleteGroup: 'Excluir Grupo',
        ifLabel: 'SE',
        operatorLabel: 'OPERADOR',
        compareToLabel: 'COMPARAR COM',
        addConditionHint: 'Adicionar uma condição',
        addConditionButton: 'Adicionar Condição',
        itemSingular: 'item',
        itemPlural: 'itens',
      },
    },
  },
};

/** Core Data plugin translations */
export const coreData = {
  nodes: {
    log: {
      label: 'Log',
      description: 'Emitir evento de observabilidade',
      config: {
        message: 'Mensagem',
        level: 'Nível',
        levelInfo: 'Info',
        levelWarn: 'Aviso',
        levelError: 'Erro',
        levelDebug: 'Debug',
      },
    },
    set_state: {
      label: 'Definir Estado',
      description: 'Definir, incrementar, decrementar, adicionar ao array ou remover uma variável de estado',
      config: {
        stateVariableSection: 'Variável de Estado',
        selectVariable: 'Selecionar variável...',
        noStateVariables: 'Nenhuma variável de estado definida',
        operationSection: 'Operação',
        operationSet: 'Definir',
        operationSetDesc: 'Substituir o valor atual',
        operationIncrement: 'Incrementar',
        operationIncrementDesc: 'Adicionar ao valor atual',
        operationDecrement: 'Decrementar',
        operationDecrementDesc: 'Subtrair do valor atual',
        operationAppend: 'Adicionar ao Array',
        operationAppendDesc: 'Adicionar um valor ao final de uma variável do tipo array',
        operationRemove: 'Remover',
        operationRemoveDesc: 'Limpar a variável de estado',
        valueSection: 'Valor',
        removedHint: 'A variável de estado será limpa quando este nó for executado.',
        appendHint: 'O valor será adicionado ao final do array. Se a variável não for um array, será convertida para um.',
      },
    },
    code: {
      label: 'Código',
      description: 'Executar um snippet JavaScript',
      config: {
        scriptSection: 'Script',
        linesBadge: 'linhas',
        openEditor: 'Abrir Editor',
        timeoutSection: 'Timeout',
        timeoutHint: 'Tempo máximo de execução (mín: 100ms)',
        timeoutMin: 'Mínimo 100ms',
        availableContext: 'Disponível: {state}, {event}, {variables}, {nodes}. Retorne um objeto para expor saída via {nodesOutput}',
        scriptEditorTitle: 'Editor de Script',
        sandboxedHint: 'Sandbox — sem módulos externos',
        defaultScript: '// Acesso: state, event, variables, nodes\n\nreturn {};',
      },
    },
  },
};

/** Core Flow Control plugin translations */
export const coreFlowControl = {
  nodes: {
    fanout: {
      label: 'Fanout',
      description: 'Fork de execução paralela',
      config: {
        branches: 'Número de Ramificações',
      },
    },
    sequence: {
      label: 'Sequência',
      description: 'Execução sequencial — executa cada etapa em ordem',
      config: {
        steps: 'Número de Etapas',
      },
    },
    merge: {
      label: 'Merge',
      description: 'Juntar ramificações paralelas',
      config: {
        branches: 'Número de Ramificações',
        strategy: 'Estratégia de Merge',
        strategyAll: 'Todos',
        strategyAny: 'Qualquer',
        strategyFirst: 'Primeiro',
      },
    },
    switch: {
      label: 'Switch',
      description: 'Rotear para caminhos com base em condições',
      config: {
        noCasesYet: 'Nenhum caso ainda',
        evaluationModeSection: 'Modo de Avaliação',
        firstMatch: 'Primeira Correspondência',
        allMatches: 'Todas as Correspondências',
        firstMatchDesc: 'Para na primeira correspondência (exclusivo)',
        allMatchesDesc: 'Ativa todas as correspondências em paralelo (inclusivo)',
        casesSection: 'Casos',
        caseLabel: 'Caso {number}',
        deleteCase: 'Excluir Caso',
        noConditionsInCase: 'Nenhuma condição neste caso',
        conditionLabel: 'Condição',
        conditionDescription: 'Comparação simples de campo',
        groupLabel: 'Grupo',
        groupDescription: 'Grupo de condições com lógica',
        addButton: 'Adicionar',
        defaultCase: 'Padrão — quando nenhum caso corresponde',
      },
    },
    subworkflow: {
      label: 'Subworkflow',
      description: 'Executar outro workflow como processo filho',
      config: {
        workflowSection: 'Workflow',
        subworkflowBadge: 'subworkflow',
        changeWorkflow: 'Alterar workflow',
        removeWorkflow: 'Remover workflow',
        selectWorkflow: 'Selecionar Workflow',
        executionModeSection: 'Modo de Execução',
        sync: 'Síncrono',
        async: 'Assíncrono',
        syncDescription: 'Aguarda o workflow filho completar. Saída disponível via nodes.(id).output.*',
        asyncDescription: 'Dispara e esquece — continua imediatamente. Sem retorno de saída.',
        executionTimeoutSection: 'Timeout de Execução',
        duration: 'Duração',
        seconds: 'Segundos',
        minutes: 'Minutos',
        hours: 'Horas',
        timeoutSyncHint: 'Se o timeout for atingido, a saída de erro é disparada.',
        timeoutAsyncHint: 'Se o timeout for atingido, a execução filha é cancelada.',
        inputMappingsSection: 'Mapeamentos de Entrada',
        inputMappingsHint: 'Passe dados do workflow pai para as variáveis do workflow filho.',
        childVariable: 'Variável filha',
        childVariablePlaceholder: 'ex. deviceId',
        valuePlaceholder: 'ex. result.status',
        removeMapping: 'Remover mapeamento',
        addInput: 'Adicionar Entrada',
        outputMappingsSection: 'Mapeamentos de Saída',
        outputMappingsHint: 'Mapear saída do workflow filho para variáveis de estado do pai.',
        outputAlsoAvailable: 'A saída também está disponível como {expression} em expressões downstream.',
        childOutputKey: 'Chave de saída filha',
        parentVariable: 'Variável pai',
        parentVariablePlaceholder: 'ex. childResult',
        addOutput: 'Adicionar Saída',
        selectPrompt: 'Selecione um workflow para configurar este nó de subworkflow',
      },
    },
    loop: {
      label: 'Loop',
      description: 'Iterar sobre uma lista',
      config: {
        sourceSection: 'Origem',
        sourceHint: 'Cada iteração expõe {loopItem} (item atual) e {loopIndex} (índice base 0) para nós downstream via a saída Body. A saída Done dispara após todas as iterações.',
      },
    },
    end: {
      label: 'Fim',
      description: 'Ponto de terminação do workflow com modo de erro opcional',
      config: {
        terminateWithErrorBanner: 'O workflow terminará com status de erro.',
        terminateSuccessBanner: 'O workflow terminará com sucesso.',
        terminationModeSection: 'Modo de Terminação',
        terminateWithError: 'Terminar com Erro',
        terminateWithErrorHint: 'Encerrar workflow em estado de erro',
        errorCodeSection: 'Código de Erro',
        errorCodePlaceholder: 'ex. VALIDATION_FAILED',
        errorCodeHint: 'Código de erro único para tratamento programático',
        errorMessageSection: 'Mensagem de Erro',
        compensationHint: 'Nós upstream com saídas de erro podem acionar lógica de compensação antes de chegar a este nó final.',
      },
    },
    goto: {
      label: 'Goto',
      description: 'Portal virtual — conectar seções do workflow sem arestas',
      config: {
        roleSection: 'Função',
        sender: 'Emissor',
        receiver: 'Receptor',
        senderDescription: 'Endpoint — o fluxo chega aqui e pula para o receptor correspondente',
        receiverDescription: 'Origem — o fluxo retoma aqui a partir de qualquer emissor correspondente',
        labelSection: 'Rótulo',
        labelPlaceholder: 'ex. ERR, P1, RetryBlock',
        labelHint: 'Crie um receptor com o mesmo rótulo para completar o portal',
        targetSenderSection: 'Emissor Alvo',
        noSendersAvailable: 'Nenhum emissor disponível — crie um Goto Emissor primeiro',
        selectSenderLabel: 'Selecione um rótulo de emissor...',
        sendBadge: 'Enviar',
        recvBadge: 'Receber',
        colorSection: 'Cor',
        senderHint: 'O fluxo chega aqui e pula para o portal receptor correspondente',
        receiverHint: 'O fluxo retoma aqui a partir de qualquer portal emissor correspondente',
        noMatchingGoto: 'Nenhum outro nó goto com rótulo',
        matchedPairsSection: 'Pares Correspondentes',
        matchedPairsHint: 'Nós Goto conectam seções do workflow sem arestas visíveis. O backend resolve pares emissor/receptor por rótulos correspondentes.',
      },
    },
  },
};

/** Core Timers plugin translations */
export const coreTimers = {
  nodes: {
    delay: {
      label: 'Atraso',
      description: 'Espera fixa (segundos a anos)',
      config: {
        duration: 'Duração',
        unit: 'Unidade',
        unitSeconds: 'Segundos',
        unitMinutes: 'Minutos',
        unitHours: 'Horas',
        unitDays: 'Dias',
        unitMonths: 'Meses (30 dias)',
        unitYears: 'Anos (365 dias)',
      },
    },
    wait_signal: {
      label: 'Aguardar Sinal',
      description: 'Aguardar sinal externo com timeout',
      config: {
        signalSection: 'Sinal',
        signalNamePlaceholder: 'Nome do sinal (ex. approval_response)',
        signalNameHint: 'Nome único do sinal a escutar',
        timingSection: 'Temporização',
        timeout: 'Timeout',
        timeoutPlaceholder: 'ex. 10m, 1h, 24h',
        maxTimeoutCycles: 'Ciclos Máx. de Timeout',
        maxTimeoutCyclesHint: 'Tentativas internas antes de ir para saída de Timeout',
        variableMappingsSection: 'Mapeamentos de Variáveis',
        variableMappingsHint: 'Mapeie campos do payload do sinal para variáveis de estado do workflow. Quando o sinal chegar, os valores mapeados serão gravados nas variáveis de estado correspondentes.',
        fromLabel: 'De',
        fromPlaceholder: 'payload.campo.caminho',
        toLabel: 'Para',
        toPlaceholder: 'state.xxx',
        removeMapping: 'Remover mapeamento',
        noMappingsHint: 'Sem mapeamentos — o sinal desbloqueará o fluxo sem gravar dados',
        addMapping: 'Adicionar Mapeamento',
        noStateVariables: 'Sem variáveis de estado. Adicione na aba Estado.',
        noSignalsDefined: 'Nenhum sinal definido. Adicione na aba Dados > Sinais.',
      },
    },
    wait_for: {
      label: 'Aguardar Condição',
      description: 'Aguardar até que uma condição de variável de estado seja atendida',
      config: {
        conditionSection: 'Condição',
        variable: 'Variável',
        selectStateVariable: 'Selecionar variável de estado...',
        noStateVariables: 'Nenhuma variável de estado definida. Adicione variáveis na aba Estado.',
        operatorSection: 'Operador',
        compareToSection: 'Comparar Com',
        source: 'Origem',
        sourceLiteral: 'Literal',
        sourceVariable: 'Variável',
        value: 'Valor',
        valuePlaceholder: 'Digite o valor...',
        timingSection: 'Temporização',
        pollingInterval: 'Intervalo de Polling',
        pollingIntervalPlaceholder: 'ex. 30s, 1m, 5m',
        timeout: 'Timeout',
        timeoutPlaceholder: 'ex. 5m, 1h, 24h',
        maxTimeoutCycles: 'Ciclos Máx. de Timeout',
        maxTimeoutCyclesHint: 'Tentativas internas antes de ir para saída de Timeout',
      },
    },
  },
};

/** Core Annotations plugin translations */
export const coreAnnotations = {
  nodes: {
    text_note: {
      label: 'Nota de Texto',
      description: 'Adicionar uma nota de texto para documentar o workflow',
      placeholder: 'Digite sua nota...',
      emptyPlaceholder: 'Clique duplo para editar...',
    },
    group_frame: {
      label: 'Quadro de Grupo',
      description: 'Container visual para organizar seções do workflow',
      defaultTitle: 'Grupo',
      config: {
        titleSection: 'Título',
        titlePlaceholder: 'ex. Bloco de Tratamento de Erros',
        descriptionSection: 'Descrição',
        descriptionPlaceholder: 'Descrição opcional...',
        colorSection: 'Cor',
        sizeSection: 'Tamanho',
        widthLabel: 'Largura (px)',
        heightLabel: 'Altura (px)',
        resizeHint: 'Você também pode redimensionar arrastando os cantos do quadro no canvas',
        infoHint: 'Quadros de grupo são containers visuais para organizar seções do workflow. Não têm impacto funcional na execução.',
      },
    },
  },
};

/** Shared strings used across multiple config components */
export const shared = {
  selectEventField: 'Selecionar Campo de Evento',
  searchFields: 'Buscar campos...',
  noFieldsAvailable: 'Nenhum campo disponível dos templates selecionados.',
  fieldSingular: 'campo',
  fieldPlural: 'campos',
  templatesSelected: 'template(s) selecionado(s)',
  change: 'Alterar',
  from: 'De:',
  addNote: 'Adicionar Nota',
  durable: 'Durável',
};
