// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
"use client";
import { useEffect, useRef, useState, useCallback } from "react";
import FileEditor from "./FileEditor";
import Ajv, { AnySchema } from "ajv";
import { useLoading } from "../Layout";
import { Alert, Button } from "antd";

const ajv = new Ajv();
let validate: any;

export default function AdvancedView() {
  const [flagData, setflagData] = useState(null);
  const [reloadData, setReloadData] = useState(false);
  const [flagDataIsSynced, setFlagDataIsSynced] = useState(true);

  const textAreaRef = useRef<HTMLTextAreaElement>(null);
  const { setIsLoading } = useLoading();

  const requestSchemas = useCallback(async (url: string): Promise<AnySchema | null> => {
    try {
      setIsLoading(true);
      const response = await fetch(url);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();
      setIsLoading(false);
      return data;
    } catch (error: any) {
      console.error("There was an error:", error.message);
    }
    return null;
  }, [setIsLoading]);

  useEffect(() => {
    const readFile = async (file_name: string) => {
      try {
        const response = await fetch(`/feature/api/read-file?${file_name}`, {
          method: "GET",
          headers: { "Content-Type": "application/json" },
        });
        const data = await response.json();
        setflagData(data);
      } catch (err: unknown) {
        window.alert(err);
        console.error(err);
      }
    };
    readFile("");
    if (reloadData) {
      setReloadData(false);
    }
  }, [reloadData]);

  useEffect(() => {
    async function loadSchema() {
      try {
        const schemas: [AnySchema | null, AnySchema | null] = await Promise.all(
          [
            requestSchemas("https://flagd.dev/schema/v0/flags.json"),
            requestSchemas("https://flagd.dev/schema/v0/targeting.json"),
          ],
        );
        if (schemas[0] && schemas[1]) {
          validate = ajv.addSchema(schemas[1]).compile(schemas[0]);
        }
      } catch (error) {
        console.warn("Error loading schemas:", error);
      }

      return null;
    }
    loadSchema();
  }, [ajv, requestSchemas]);


  function parseJSON(): string | null {
    try {
      if (textAreaRef.current) {
        const data = JSON.parse(textAreaRef.current.value);
        return data;
      }
    } catch (objError) {
      window.alert("Error parsing JSON");
      if (objError instanceof SyntaxError) {
        console.error(objError.name);
      } else {
        console.error(objError);
      }
    }
    return null;
  }

  async function saveUpdate(flagData: string): Promise<void> {
    try {
      setIsLoading(true);
      const response = await fetch("/feature/api/write-to-file", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ data: flagData }),
      });
      await response.json();
      setIsLoading(false);
      setReloadData(true);
      setFlagDataIsSynced(true);
    } catch (err: unknown) {
      setIsLoading(false);
      window.alert(err);
      console.error(err);
    }
  }

  function update() {
    const data = parseJSON();
    if (data === null) return;
    validate(data) ? saveUpdate(data) : window.alert("Schema check failed");
  }

  const handleTextAreaChange = () => {
    try {
      if (textAreaRef.current) {
        const textAreaContent = JSON.parse(textAreaRef.current.value);
        setFlagDataIsSynced(
          JSON.stringify(textAreaContent) === JSON.stringify(flagData),
        );
      }
    } catch (error) {
      console.error("Invalid JSON in textarea", error);
      setFlagDataIsSynced(false);
    }
  };

  return (
    <>
      {flagData && (
        <div>
          <div className="p-2 pl-8 text-gray-300 shadow-md">
            <div className="mb-8 flex flex-auto items-center gap-2">
              <Button className='mr-4' type="primary" size='large' onClick={update}>保存</Button>
              {!flagDataIsSynced && <Alert message="有未保存的项目" type="warning" showIcon/>}
            </div>
          </div>
          <FileEditor
            flagConfig={flagData}
            textAreaRef={textAreaRef as React.RefObject<HTMLTextAreaElement>}
            handleTextAreaChange={handleTextAreaChange}
          />
        </div>
      )}
    </>
  );
}
