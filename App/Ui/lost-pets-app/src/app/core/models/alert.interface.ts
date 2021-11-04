export interface IAlert {
    type: alertType
    message: string
    duration?: number // in seconds
}

export enum alertType {
    info = 'info', // default behavior for clr-alerts
    warning = 'warning',
    success = 'success',
    danger = 'danger'
}
