// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
"use client";
import React, { useEffect } from "react";
import AdvancedView from "@/components/advanced/AdvancedView";
import { useRouter } from 'next/navigation';
import { checkAuth } from '../../lib/checkAuth';

const Advanced = () => {

  const router = useRouter();
  useEffect(() => {
    if (!checkAuth()) {
      router.push('/login');
    }
  }, [router]);
  return <AdvancedView />;
};

export default Advanced;
