// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

import { v4 } from 'uuid';

interface ISession {
  userId: string;
  currencyCode: string;
  sid: string;
}

const sessionKey = 'session';
const sIdKey = 'sid'
const defaultSession = {
  userId: v4(),
  currencyCode: 'USD',
  sid: '',
};

const SessionGateway = () => ({
  getSession(): ISession {
    if (typeof window === 'undefined') return defaultSession;
    const sessionString = localStorage.getItem(sessionKey);
    if (!sessionString) localStorage.setItem(sessionKey, JSON.stringify(defaultSession));
    return JSON.parse(sessionString || JSON.stringify(defaultSession)) as ISession;
  },
  setSessionValue<K extends keyof ISession>(key: K, value: ISession[K]) {
    const session = this.getSession();

    localStorage.setItem(sessionKey, JSON.stringify({ ...session, [key]: value }));
  },
  removeSessionValue<K extends keyof ISession>(key: K) {
      const session = this.getSession();
      const { [key]: _, ...rest } = session;
      localStorage.setItem(sessionKey, JSON.stringify(rest));
   }
});

export default SessionGateway();
