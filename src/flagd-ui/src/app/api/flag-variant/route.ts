// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
import { NextResponse } from "next/server";
import fs from "fs";
import * as fsPromises from "fs/promises";
import path from "path";
import { FlagConfig } from "@/utils/types";
import { Variants } from "antd/es/config-provider";

export async function GET(request: Request) {
  const { searchParams } = new URL(request.url);
  const file_name = searchParams.get("file_name") || "demo.flagd.json";
  const flag_id = searchParams.get("flag_id") || undefined;
  if (!flag_id) {
    return NextResponse.json({ error: "flag_id param is empty." }, { status: 403 });
  }
  let data:any
  try {
    const filePath = path.join(process.cwd(), "data", file_name);
    const fileContents = fs.readFileSync(filePath, "utf8");
    data = JSON.parse(fileContents);
  } catch (error) {
    console.error("Error reading file:", error);
    return NextResponse.json({ error: "Failed to read data file" }, { status: 500 });
  }
  try {
    const flags = data?.flags ?? Object()
    const ids = Object.keys(flags)
    for (const id of ids) {
      if (id == flag_id) {
        const flagConfig: FlagConfig = flags[id];
        return NextResponse.json({
          id: id,
          variants: flagConfig.variants,
          defaultVariant: flagConfig.defaultVariant
        });
      }
    }
    return NextResponse.json({
      id: flag_id,
      variants: {},
      defaultVariant: undefined
    });
  } catch (error) {
    console.error("Error get flag: ", error)
    return NextResponse.json({ 
      id: flag_id,
      success: false,
      error: "Failed to get flag, " + error
    }, { status: 500 });
  }
}

export async function POST(request: Request) {
  const { searchParams } = new URL(request.url);
  const file_name = searchParams.get("file_name") || "demo.flagd.json";
  const flag_id = searchParams.get("flag_id") || undefined;
  const new_var = searchParams.get("default_variant") || undefined;
  if (!flag_id) {
    return NextResponse.json({ error: "flag_id param is empty." }, { status: 403 });
  }
  if (!new_var) {
    return NextResponse.json({ error: "default_variant param is empty." }, { status: 403 });
  }
  let data = Object()
  try {
    const filePath = path.join(process.cwd(), "data", file_name);
    const fileContents = fs.readFileSync(filePath, "utf8");
    data = JSON.parse(fileContents);
  } catch (error) {
    console.error("Error reading file:", error);
    return NextResponse.json({ error: "Failed to read data file" }, { status: 500 });
  }
  if (!data.flags) {
    data.flags = Object()
  }
  try {
    let flags = data.flags
    const ids = Object.keys(flags)
    if (ids.indexOf(flag_id) === -1) {
      return NextResponse.json({
        id: flag_id,
        success: false,
        error: "Failed to find '"+flag_id+"' in [" + ids.toString() + "]"
      }, { status: 403 });
    }
    const flagConfig: FlagConfig = flags[flag_id];
    const keys = Object.keys(flagConfig?.variants ?? {})
    if (keys.indexOf(new_var) === -1) {
      return NextResponse.json(
        {
          id: flag_id,
          success: false,
          error: "Failed to find '" + new_var + "' in [" + keys.toString()+ "]",
        },
        { status: 403 },
      );
    }
    flags[flag_id]["defaultVariant"] = new_var
    data.flags = flags
    const filePath = path.join(process.cwd(), "data", file_name);
    
    await fsPromises.writeFile(filePath, JSON.stringify(data, null, 2), "utf8");
    
    return NextResponse.json({
      id: flag_id,
      success: true
    });
  } catch (error) {
    console.error("Error update flag: ", error)
    return NextResponse.json({ 
      id: flag_id,
      success: false,
      error: "Failed to update flag, " + error
    }, { status: 500 });
  }
}
