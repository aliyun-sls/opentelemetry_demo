// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
"use client";
import React, { ChangeEvent, useCallback, useEffect } from "react";
import { useState } from "react";
import DefaultVariantSelect from "./DefaultVariantSelect";
import { FlagConfig } from "@/utils/types";
import { Button, Tag, Typography } from "antd";

type FeatureFlagProps = {
  flagId: string;
  flagConfig: FlagConfig;
  updateFlagData: (fladId: string, selectedVariant: string) => void;
};

function hashCode(str: string) {
  // Compute a simple hash value from the string
  let hash = 0;
  if (!str) return 0;
  for (let i = 0; i < str.length; i++) {
    let char = str.charCodeAt(i);
    hash = ((hash << 5) - hash) + char;
    hash |= 0; // Convert to 32-bit integer
  }
  return hash;
}

const COLOR = [
  "magenta",
  "red",
  "volcano",
  "orange",
  "gold",
  "lime",
  "green",
  "cyan",
  "blue",
  "geekblue",
  "purple",
];

export function getColorOnText(text: string, range = COLOR.length) {
  let hashValue = hashCode(text);
  // Use modulo to ensure the number is within a specific range
  let n = Math.abs(hashValue) % range;
  return COLOR[n];
}

function FeatureFlag({ flagId, flagConfig, updateFlagData }: FeatureFlagProps) {
  const [selectedVariant, setSelectedVariant] = useState<string>("");

  useEffect(() => {
    setSelectedVariant(flagConfig.defaultVariant);
  }, [flagConfig.defaultVariant]);

  const handleVariantChange = useCallback(
    (event: ChangeEvent<HTMLSelectElement>) => {
      setSelectedVariant(event.target.value);
      updateFlagData(flagId, event.target.value);
    },
    [flagId, updateFlagData],
  );

  return (
    <div className="mb-4 flex flex-auto flex-col justify-between rounded-md bg-gray-800 p-6 text-gray-300 shadow-md">
      <div>
        <div className="flex items-center justify-between">
          <div className="mb-1 text-lg font-semibold">{`${flagId}`}</div>
          <Typography.Text copyable={{ text: `${flagId}` }} />
        </div>
        <hr  className="mb-3 text-gray-600"/>
        <div className="mb-2 flex flex-wrap gap-0.5">
          {flagConfig.tags && flagConfig.tags.map((tag) => (
            <Tag key={tag} bordered={false} color={getColorOnText(tag)}>{tag}</Tag>
          ))}
        </div>
        <p className="mb-4 text-sm">{`${flagConfig.description}`}</p>
      </div>
      <div>
        <div className="flex items-center justify-between">
          <DefaultVariantSelect
            flagConfig={flagConfig}
            selectedVariant={selectedVariant}
            handleVariantChange={handleVariantChange}
          />
        </div>
      </div>
    </div>
  );
}

export default FeatureFlag;
