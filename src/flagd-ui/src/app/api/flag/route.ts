// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
import { NextResponse } from "next/server";
import fs from "fs";
import * as fsPromises from "fs/promises";
import path from "path";
import { FlagConfig } from "@/utils/types";

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
          config: flagConfig
        });
      }
    }
    return NextResponse.json({
      id: flag_id,
      config: {}
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

export async function PUT(request: Request) {
  const { searchParams } = new URL(request.url);
  const file_name = searchParams.get("file_name") || "demo.flagd.json";
  const flag_id = searchParams.get("flag_id") || undefined;
  const config = await request.json();
  if (!flag_id) {
    return NextResponse.json({ error: "flag_id param is empty." }, { status: 403 });
  }
  if (!config) {
    return NextResponse.json({ error: "request body is empty. cannot get config." }, { status: 403 });
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
    const flags = data.flags
    flags[flag_id] = config
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


export async function DELETE(request: Request) {
  const { searchParams } = new URL(request.url);
  const file_name = searchParams.get("file_name") || "demo.flagd.json";
  const flag_id = searchParams.get("flag_id") || undefined;
  if (!flag_id) {
    return NextResponse.json({ error: "flag_id param is empty." }, { status: 403 });
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
    const flags = data.flags
    const ids = Object.keys(flags)
    let new_flags = Object()
    for (const id of ids) {
      if (id != flag_id) {
        new_flags[id] = flags[id];
      }
    }
    data.flags = new_flags
    const filePath = path.join(process.cwd(), "data", file_name);
    
    await fsPromises.writeFile(filePath, JSON.stringify(data, null, 2), "utf8");
    
    return NextResponse.json({
      id: flag_id,
      success: true
    });
  } catch (error) {
    console.error("Error delete flag: ", error)
    return NextResponse.json({ 
      id: flag_id,
      success: false,
      error: "Failed to delete flag, " + error
    }, { status: 500 });
  }
}