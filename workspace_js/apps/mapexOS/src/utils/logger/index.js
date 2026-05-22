import dayjs from 'dayjs'

const DATE_BR = () => dayjs().format('DD-MM-YYYY HH:mm:ss')

const FORMATER = (type, message) => {
  if (process.env.DEV) {
    const dateToPrint = {
      type,
      date: DATE_BR(),
      message
    }

    console[type](dateToPrint)
  }
}

export const log = message => FORMATER('log', message)
export const error = message => FORMATER('error', message)
export const warn = message => FORMATER('warn', message)
export const info = message => FORMATER('info', message)
