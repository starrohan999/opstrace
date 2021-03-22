/**
 * Copyright 2020 Opstrace, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
import { ReactNode } from "react";

export type PickerOption<T = {}> = {
  id: string;
  text: string;
  data: T;
};

export type PickerListProps = {
  selectedIndex: number;
  options: PickerOption[];
  onSelect: (selectedOption: PickerOption) => void;
  secondaryAction?: (option: PickerOption) => ReactNode;
};

export type PickerProvider<T = {}> = {
  title?: string;
  // Disable's the filtering of options
  disableFilter?: boolean;
  // Disable input
  disableInput?: boolean;
  // Validate input as user is typing
  textValidator?: RegExp;
  textValidationFailedMessage?: string;
  activationPrefix: string;
  onSelected: (option: PickerOption<T>, inputText?: string) => void;
  options: PickerOption<T>[];
  secondaryAction?: (option: PickerOption<T>) => ReactNode;
};

export type PickerApi = {
  register: (provider: PickerProvider) => void;
  unregister: (provider: PickerProvider) => void;
  setText: (text: string) => void;
};

export type PickerState = {
  activeProviderIndex: number;
  text: null | string;
  providers: PickerProvider[];
};
