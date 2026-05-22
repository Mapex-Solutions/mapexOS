/**
 * English (en-US) translations for core workflow plugins.
 * Organized by plugin ID → node short name → translation keys.
 *
 * Namespace at runtime: wf.{pluginId}.nodes.{nodeShortName}.*
 */

/** Core Triggers plugin translations */
export const coreTriggers = {
  nodes: {
    start: {
      label: 'Start',
      description: 'Workflow entry point',
    },
    trigger_event: {
      label: 'Trigger Event',
      description: 'Entry point — triggered by a configured trigger from the platform',
      config: {
        triggerSection: 'Trigger',
        changeTrigger: 'Change trigger',
        removeTrigger: 'Remove trigger',
        selectTrigger: 'Select Trigger',
        variablesSection: 'Variables',
        selectPrompt: 'Select a trigger to configure this event node',
      },
    },
  },
};

/** Core Logic plugin translations */
export const coreLogic = {
  nodes: {
    condition: {
      label: 'Conditions',
      description: 'Evaluate conditions and branch based on result',
      config: {
        noConditionsYet: 'No conditions yet',
        conditionLabel: 'Condition',
        conditionDescription: 'Simple field comparison',
        groupLabel: 'Group',
        groupDescription: 'Group of conditions with logic',
        addButton: 'Add',
        rename: 'Rename',
        deleteItem: 'Delete',
        deleteGroup: 'Delete Group',
        ifLabel: 'IF',
        operatorLabel: 'OPERATOR',
        compareToLabel: 'COMPARE TO',
        addConditionHint: 'Add a condition',
        addConditionButton: 'Add Condition',
        itemSingular: 'item',
        itemPlural: 'items',
      },
    },
  },
};

/** Core Data plugin translations */
export const coreData = {
  nodes: {
    log: {
      label: 'Log',
      description: 'Emit observability event',
      config: {
        message: 'Message',
        level: 'Level',
        levelInfo: 'Info',
        levelWarn: 'Warn',
        levelError: 'Error',
        levelDebug: 'Debug',
      },
    },
    set_state: {
      label: 'Set State',
      description: 'Set, increment, decrement, append or remove a state variable',
      config: {
        stateVariableSection: 'State Variable',
        selectVariable: 'Select variable...',
        noStateVariables: 'No state variables defined',
        operationSection: 'Operation',
        operationSet: 'Set',
        operationSetDesc: 'Replace the current value',
        operationIncrement: 'Increment',
        operationIncrementDesc: 'Add to the current value',
        operationDecrement: 'Decrement',
        operationDecrementDesc: 'Subtract from the current value',
        operationAppend: 'Append to Array',
        operationAppendDesc: 'Add a value to the end of an array variable',
        operationRemove: 'Remove',
        operationRemoveDesc: 'Clear the state variable',
        valueSection: 'Value',
        removedHint: 'The state variable will be cleared when this node executes.',
        appendHint: 'The value will be appended to the end of the array. If the variable is not an array, it will be converted to one.',
      },
    },
    code: {
      label: 'Code',
      description: 'Execute a JavaScript snippet',
      config: {
        scriptSection: 'Script',
        linesBadge: 'lines',
        openEditor: 'Open Editor',
        timeoutSection: 'Timeout',
        timeoutHint: 'Maximum execution time (min: 100ms)',
        timeoutMin: 'Minimum 100ms',
        availableContext: 'Available: {state}, {event}, {variables}, {nodes}. Return an object to expose output via {nodesOutput}',
        scriptEditorTitle: 'Script Editor',
        sandboxedHint: 'Sandboxed — no external modules',
        defaultScript: '// Access: state, event, variables, nodes\n\nreturn {};',
      },
    },
  },
};

/** Core Flow Control plugin translations */
export const coreFlowControl = {
  nodes: {
    fanout: {
      label: 'Fanout',
      description: 'Parallel execution fork',
      config: {
        branches: 'Number of Branches',
      },
    },
    sequence: {
      label: 'Sequence',
      description: 'Sequential execution — runs each step in order',
      config: {
        steps: 'Number of Steps',
      },
    },
    merge: {
      label: 'Merge',
      description: 'Join parallel branches',
      config: {
        branches: 'Number of Branches',
        strategy: 'Merge Strategy',
        strategyAll: 'All',
        strategyAny: 'Any',
        strategyFirst: 'First',
      },
    },
    switch: {
      label: 'Switch',
      description: 'Route to paths based on matching conditions',
      config: {
        noCasesYet: 'No cases yet',
        evaluationModeSection: 'Evaluation Mode',
        firstMatch: 'First Match',
        allMatches: 'All Matches',
        firstMatchDesc: 'Stops at the first matching case (exclusive)',
        allMatchesDesc: 'Activates all matching cases in parallel (inclusive)',
        casesSection: 'Cases',
        caseLabel: 'Case {number}',
        deleteCase: 'Delete Case',
        noConditionsInCase: 'No conditions in this case',
        conditionLabel: 'Condition',
        conditionDescription: 'Simple field comparison',
        groupLabel: 'Group',
        groupDescription: 'Group of conditions with logic',
        addButton: 'Add',
        defaultCase: 'Default — when no case matches',
      },
    },
    subworkflow: {
      label: 'Subworkflow',
      description: 'Execute another workflow as a child process',
      config: {
        workflowSection: 'Workflow',
        subworkflowBadge: 'subworkflow',
        changeWorkflow: 'Change workflow',
        removeWorkflow: 'Remove workflow',
        selectWorkflow: 'Select Workflow',
        executionModeSection: 'Execution Mode',
        sync: 'Sync',
        async: 'Async',
        syncDescription: 'Wait for child workflow to complete. Output available via nodes.(id).output.*',
        asyncDescription: 'Fire-and-forget — continue immediately. No output returned.',
        executionTimeoutSection: 'Execution Timeout',
        duration: 'Duration',
        seconds: 'Seconds',
        minutes: 'Minutes',
        hours: 'Hours',
        timeoutSyncHint: 'If timeout is reached, the error output is triggered.',
        timeoutAsyncHint: 'If timeout is reached, the child execution is cancelled.',
        inputMappingsSection: 'Input Mappings',
        inputMappingsHint: 'Pass data from the parent workflow into the child workflow\'s variables.',
        childVariable: 'Child variable',
        childVariablePlaceholder: 'e.g. deviceId',
        valuePlaceholder: 'e.g. result.status',
        removeMapping: 'Remove mapping',
        addInput: 'Add Input',
        outputMappingsSection: 'Output Mappings',
        outputMappingsHint: 'Map child workflow output to parent state variables.',
        outputAlsoAvailable: 'Output is also available as {expression} in downstream expressions.',
        childOutputKey: 'Child output key',
        parentVariable: 'Parent variable',
        parentVariablePlaceholder: 'e.g. childResult',
        addOutput: 'Add Output',
        selectPrompt: 'Select a workflow to configure this subworkflow node',
      },
    },
    loop: {
      label: 'Loop',
      description: 'Iterate over a list',
      config: {
        sourceSection: 'Source',
        sourceHint: 'Each iteration exposes {loopItem} (current item) and {loopIndex} (0-based index) to downstream nodes via the Body output. The Done output fires after all iterations complete.',
      },
    },
    end: {
      label: 'End',
      description: 'Workflow termination point with optional error mode',
      config: {
        terminateWithErrorBanner: 'Workflow will terminate with an error status.',
        terminateSuccessBanner: 'Workflow will terminate successfully.',
        terminationModeSection: 'Termination Mode',
        terminateWithError: 'Terminate with Error',
        terminateWithErrorHint: 'End workflow in an error state',
        errorCodeSection: 'Error Code',
        errorCodePlaceholder: 'e.g. VALIDATION_FAILED',
        errorCodeHint: 'Unique error code for programmatic handling',
        errorMessageSection: 'Error Message',
        compensationHint: 'Upstream nodes with error outputs can trigger compensation logic before reaching this end node.',
      },
    },
    goto: {
      label: 'Goto',
      description: 'Virtual portal — connect workflow sections without edges',
      config: {
        roleSection: 'Role',
        sender: 'Sender',
        receiver: 'Receiver',
        senderDescription: 'Endpoint — flow arrives here and jumps to the matching receiver',
        receiverDescription: 'Source — flow resumes here from any matching sender',
        labelSection: 'Label',
        labelPlaceholder: 'e.g. ERR, P1, RetryBlock',
        labelHint: 'Create a receiver with the same label to complete the portal',
        targetSenderSection: 'Target Sender',
        noSendersAvailable: 'No senders available — create a Goto Sender first',
        selectSenderLabel: 'Select a sender label...',
        sendBadge: 'Send',
        recvBadge: 'Recv',
        colorSection: 'Color',
        senderHint: 'Flow arrives here and jumps to the matching receiver portal',
        receiverHint: 'Flow resumes here from any matching sender portal',
        noMatchingGoto: 'No other goto nodes with label',
        matchedPairsSection: 'Matched Pairs',
        matchedPairsHint: 'Goto nodes connect workflow sections without visible edges. The backend resolves sender/receiver pairs by matching labels.',
      },
    },
  },
};

/** Core Timers plugin translations */
export const coreTimers = {
  nodes: {
    delay: {
      label: 'Delay',
      description: 'Fixed wait (seconds to years)',
      config: {
        duration: 'Duration',
        unit: 'Unit',
        unitSeconds: 'Seconds',
        unitMinutes: 'Minutes',
        unitHours: 'Hours',
        unitDays: 'Days',
        unitMonths: 'Months (30 days)',
        unitYears: 'Years (365 days)',
      },
    },
    wait_signal: {
      label: 'Wait Signal',
      description: 'Wait for external signal with timeout',
      config: {
        signalSection: 'Signal',
        signalNamePlaceholder: 'Signal name (e.g., approval_response)',
        signalNameHint: 'Unique signal name to listen for',
        timingSection: 'Timing',
        timeout: 'Timeout',
        timeoutPlaceholder: 'e.g., 10m, 1h, 24h',
        maxTimeoutCycles: 'Max Timeout Cycles',
        maxTimeoutCyclesHint: 'Internal retries before going to Timeout output',
        variableMappingsSection: 'Variable Mappings',
        variableMappingsHint: 'Map signal payload fields to workflow state variables. When the signal arrives, the mapped payload values will be written into the corresponding state variables.',
        fromLabel: 'From',
        fromPlaceholder: 'payload.field.path',
        toLabel: 'To',
        toPlaceholder: 'state.xxx',
        removeMapping: 'Remove mapping',
        noMappingsHint: 'No mappings — signal will unblock the flow without writing data',
        addMapping: 'Add Mapping',
        noStateVariables: 'No state variables. Add them in the State tab.',
        noSignalsDefined: 'No signals defined. Add them in the Data > Signals tab.',
      },
    },
    wait_for: {
      label: 'Wait For',
      description: 'Wait until a state variable condition is met',
      config: {
        conditionSection: 'Condition',
        variable: 'Variable',
        selectStateVariable: 'Select state variable...',
        noStateVariables: 'No state variables defined. Add variables in the State tab.',
        operatorSection: 'Operator',
        compareToSection: 'Compare To',
        source: 'Source',
        sourceLiteral: 'Literal',
        sourceVariable: 'Variable',
        value: 'Value',
        valuePlaceholder: 'Enter value...',
        timingSection: 'Timing',
        pollingInterval: 'Polling Interval',
        pollingIntervalPlaceholder: 'e.g., 30s, 1m, 5m',
        timeout: 'Timeout',
        timeoutPlaceholder: 'e.g., 5m, 1h, 24h',
        maxTimeoutCycles: 'Max Timeout Cycles',
        maxTimeoutCyclesHint: 'Internal retries before going to Timeout output',
      },
    },
  },
};

/** Core Annotations plugin translations */
export const coreAnnotations = {
  nodes: {
    text_note: {
      label: 'Text Note',
      description: 'Add a text note to document your workflow',
      placeholder: 'Type your note...',
      emptyPlaceholder: 'Double-click to edit...',
    },
    group_frame: {
      label: 'Group Frame',
      description: 'Visual container to organize workflow sections',
      defaultTitle: 'Group',
      config: {
        titleSection: 'Title',
        titlePlaceholder: 'e.g. Error Handling Block',
        descriptionSection: 'Description',
        descriptionPlaceholder: 'Optional description...',
        colorSection: 'Color',
        sizeSection: 'Size',
        widthLabel: 'Width (px)',
        heightLabel: 'Height (px)',
        resizeHint: 'You can also resize by dragging the frame corners on the canvas',
        infoHint: 'Group frames are visual containers for organizing workflow sections. They have no functional impact on execution.',
      },
    },
  },
};

/** Shared strings used across multiple config components */
export const shared = {
  selectEventField: 'Select Event Field',
  searchFields: 'Search fields...',
  noFieldsAvailable: 'No fields available from selected templates.',
  fieldSingular: 'field',
  fieldPlural: 'fields',
  templatesSelected: 'template(s) selected',
  change: 'Change',
  from: 'From:',
  addNote: 'Add Note',
  durable: 'Durable',
};
