-- Create workflow definition table
CREATE TABLE IF NOT EXISTS blc_workflow (
    id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL,
    version VARCHAR(50) NOT NULL DEFAULT '1.0.0',
    is_active BOOLEAN NOT NULL DEFAULT true,
    activities JSONB NOT NULL,
    transitions JSONB NOT NULL,
    start_activity_id VARCHAR(100) NOT NULL,
    end_activity_ids JSONB NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_workflow_name ON blc_workflow(name);
CREATE INDEX idx_workflow_type ON blc_workflow(type, is_active);
CREATE INDEX idx_workflow_active ON blc_workflow(is_active);

-- Create workflow execution table
CREATE TABLE IF NOT EXISTS blc_workflow_execution (
    id VARCHAR(100) PRIMARY KEY,
    workflow_id VARCHAR(100) NOT NULL,
    workflow_version VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    context JSONB NOT NULL DEFAULT '{}',
    input_data JSONB NOT NULL DEFAULT '{}',
    output_data JSONB NOT NULL DEFAULT '{}',
    current_activity_id VARCHAR(100),
    activity_history JSONB NOT NULL DEFAULT '[]',
    error_message TEXT,
    retry_count INT NOT NULL DEFAULT 0,
    started_by VARCHAR(255) NOT NULL,
    started_at TIMESTAMP WITH TIME ZONE NOT NULL,
    completed_at TIMESTAMP WITH TIME ZONE,
    last_heartbeat TIMESTAMP WITH TIME ZONE NOT NULL,
    entity_type VARCHAR(100),
    entity_id VARCHAR(255),
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_workflow_execution_workflow ON blc_workflow_execution(workflow_id);
CREATE INDEX idx_workflow_execution_status ON blc_workflow_execution(status);
CREATE INDEX idx_workflow_execution_entity ON blc_workflow_execution(entity_type, entity_id);
CREATE INDEX idx_workflow_execution_started_by ON blc_workflow_execution(started_by);
CREATE INDEX idx_workflow_execution_active ON blc_workflow_execution(status) WHERE status IN ('RUNNING', 'SUSPENDED');
CREATE INDEX idx_workflow_execution_heartbeat ON blc_workflow_execution(last_heartbeat) WHERE status = 'RUNNING';

-- Add comments
COMMENT ON TABLE blc_workflow IS 'Workflow definitions for orchestrating business processes';
COMMENT ON TABLE blc_workflow_execution IS 'Runtime instances of workflow executions';

COMMENT ON COLUMN blc_workflow.type IS 'Workflow type: CHECKOUT, ORDER_FULFILLMENT, PAYMENT_PROCESSING, RETURN_PROCESS, CUSTOM';
COMMENT ON COLUMN blc_workflow.activities IS 'JSON array of workflow activities (tasks, decisions, parallel, wait, sub-workflow, script)';
COMMENT ON COLUMN blc_workflow.transitions IS 'JSON array of state transitions between activities';
COMMENT ON COLUMN blc_workflow.start_activity_id IS 'ID of the starting activity';
COMMENT ON COLUMN blc_workflow.end_activity_ids IS 'JSON array of ending activity IDs';

COMMENT ON COLUMN blc_workflow_execution.status IS 'Execution status: PENDING, RUNNING, COMPLETED, FAILED, CANCELLED, SUSPENDED';
COMMENT ON COLUMN blc_workflow_execution.context IS 'Workflow-level variables and data passed between activities';
COMMENT ON COLUMN blc_workflow_execution.activity_history IS 'JSON array of completed activity executions';
COMMENT ON COLUMN blc_workflow_execution.last_heartbeat IS 'Last heartbeat timestamp for monitoring stale executions';
COMMENT ON COLUMN blc_workflow_execution.entity_type IS 'Type of associated entity (e.g., checkout_session, order)';
COMMENT ON COLUMN blc_workflow_execution.entity_id IS 'ID of associated entity';
