// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
"use client";
import React, { useEffect } from "react";
import BasicView from "../components/basic/BasicView";
import { useRouter } from 'next/navigation';
import { checkAuth } from '../lib/checkAuth';

const Home = () => {
  const router = useRouter();

  useEffect(() => {
    if (!checkAuth()) {
      router.push('/login');
    }
  }, [router]);

  return (
    <div className="app">
      <BasicView />
    </div>
  );
};

export default Home;
