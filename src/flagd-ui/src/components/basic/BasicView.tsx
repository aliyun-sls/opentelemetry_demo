// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
"use client";
import { Alert, Button, Collapse, Select, Tag } from 'antd';
import React, { useEffect, useMemo, useState } from "react";
import FeatureFlag, { getColorOnText } from "./FeatureFlag";
import { ConfigFile, FlagConfig, FlagGroup } from "@/utils/types";
import { sleep, useLoading } from "../Layout";
import type { SelectProps } from 'antd';

const BasicView = () => {
  const [flagData, setFlagData] = useState<ConfigFile | null>(null);
  const [originalFlagData, setOriginalFlagData] = useState<ConfigFile | null>(
    null,
  );
  const [reloadData, setReloadData] = useState(false);
  const [selectedTags, setSelectedTags] = useState<string[]>([]);

  const { setIsLoading } = useLoading();

  useEffect(() => {
    const readFile = async () => {
      try {
        const response = await fetch("/feature/api/read-file", {
          method: "GET",
          headers: { "Content-Type": "application/json" },
        });
        const data = await response.json();
        setFlagData(data);
        setOriginalFlagData(data);
      } catch (err: unknown) {
        window.alert(err);
        console.error(err);
      }
    };
    readFile();

    if (reloadData) {
      setReloadData(false);
    }
  }, [reloadData]);


  const flagGroup = useMemo(() => {
    if (!flagData) return new Map<string, FlagGroup>();
    const categoryMaps = new Map<string, FlagGroup>();
    Object.keys(flagData?.flags ?? {}).map((flagId) => {
      const flagConfig: FlagConfig = flagData.flags[flagId];
      const category = flagConfig.category;
      const config_tags = flagConfig.tags;
      if (selectedTags.length !== 0) {
        if (!config_tags) {
          return;
        }
        for (const tag of selectedTags) {
          if (!config_tags.includes(tag)) {
            return;
          }
        }
      }
      if (categoryMaps.has(category)) {
        const groups = categoryMaps.get(category);
        if (groups) {
          groups.flags[flagId] = flagConfig;
        }
      } else {
        categoryMaps.set(category, {
          category,
          flags: { [flagId]: flagConfig },
        });
      }
    });
    return categoryMaps;
  }, [flagData, selectedTags]);

  const tags = useMemo(() => {
    if (!flagData) return [];
    const tags_set = new Set(Object.keys(flagData?.flags ?? {}).map((flagId) => flagData.flags[flagId].tags).flat())
    return Array.from(tags_set).map((tag) => ({
      value: tag,
      label: tag,
    })).filter(item => item.value !== undefined)
  }, [flagData]);


  const updateFlagData = (flagId: string, selectedVariant: string) => {
    setFlagData((prevFlagData) => {
      if (!prevFlagData) return prevFlagData;

      return {
        ...prevFlagData,
        flags: {
          ...prevFlagData.flags,
          [flagId]: {
            ...prevFlagData.flags[flagId],
            defaultVariant: selectedVariant,
          },
        },
      };
    });
  };

  const save = async () => {
    try {
      setIsLoading(true);
      await fetch("/feature/api/write-to-file", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ data: flagData }),
      });
      // await sleep(100);
      setIsLoading(false);
      setReloadData(true);
    } catch (err: unknown) {
      window.alert(err);
      console.error(err);
    }
  };

  const flagDataIsSynced =
    JSON.stringify(flagData) === JSON.stringify(originalFlagData);


  type TagRender = SelectProps['tagRender'];
  const tagRender: TagRender = (props) => {
    const { label, value, closable, onClose } = props;
    const onPreventMouseDown = (event: React.MouseEvent<HTMLSpanElement>) => {
      event.preventDefault();
      event.stopPropagation();
    };
    return (
      <Tag
        bordered={false}
        color={getColorOnText(value)}
        onMouseDown={onPreventMouseDown}
        closable={closable}
        onClose={onClose}
        style={{ marginInlineEnd: 4 }}
      >
        {label}
      </Tag>
    );
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex flex-auto mb-8 items-center justify-between">
        <div className="flex flex-auto">
          <Button className='mr-4' type="primary" size='large' onClick={save}>保存</Button>
          {!flagDataIsSynced && <Alert message="有未保存的项目" type="warning" showIcon/>}
        </div>
        <Select
          mode="multiple"
          allowClear
          tagRender={tagRender}
          style={{ width: '50%', maxWidth: '50%', display: 'flex'}}
          placeholder="请选择tags"
          defaultValue={[]}
          onChange={setSelectedTags}
          options={tags}
        />
      </div>
      {flagGroup &&
        Array.from(flagGroup.entries()).map(([key, flagGroup]) => (
          <div key={key}>
            <h2 className="mb-4 text-2xl font-bold">{flagGroup.category ?? '未分类'}</h2>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
              {Object.keys(flagGroup.flags).map((flagId) => {
                const flagConfig: FlagConfig = flagGroup.flags[flagId];
                return (
                  <FeatureFlag
                    flagId={flagId}
                    key={flagId}
                    flagConfig={flagConfig}
                    updateFlagData={updateFlagData}
                  />
                );
              })}
            </div>
          </div>
        ))}
    </div>
  );
};
export default BasicView;
