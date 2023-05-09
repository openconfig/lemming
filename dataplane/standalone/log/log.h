// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

#ifndef DATAPLANE_STANDALONE_LOG_LOG_H_
#define DATAPLANE_STANDALONE_LOG_LOG_H_

#include <fstream>

// TODO(dgrau): Replace with a logging library.
extern std::ofstream logger;
#define LUCIUS_LOG_FUNC()                                            \
  logger << "Line: " << __LINE__ << " Func: " << __PRETTY_FUNCTION__ \
         << std::endl;                                               \
  logger.flush()

#define LOG(msg)                                               \
  logger << "Line: " << __LINE__ << " : " << msg << std::endl; \
  logger.flush()

#endif  // DATAPLANE_STANDALONE_LOG_LOG_H_
