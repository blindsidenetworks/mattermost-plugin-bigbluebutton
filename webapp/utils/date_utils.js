/*
Copyright 2018 Blindside Networks

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
export function formatDate(date, useMilitaryTime = false) {
  const monthNames = [
    'Jan',
    'Feb',
    'Mar',
    'Apr',
    'May',
    'Jun',
    'Jul',
    'Aug',
    'Sep',
    'Oct',
    'Nov',
    'Dec'
  ];

  const day = date.getDate();
  const monthIndex = date.getMonth();
  let hours = date.getHours();
  let minutes = date.getMinutes();

  let ampm = '';
  if (!useMilitaryTime) {
    ampm = ' AM';
    if (hours >= 12) {
      ampm = ' PM';
    }

    hours %= 12;
    if (!hours) {
      hours = 12;
    }
  }

  if (minutes < 10) {
    minutes = '0' + minutes;
  }

  return monthNames[monthIndex] + ' ' + day + ' at ' + hours + ':' + minutes + ampm;
}
