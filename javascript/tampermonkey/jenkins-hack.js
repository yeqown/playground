// ==UserScript==
// @name         Jenkins Hack
// @namespace    http://tampermonkey.net/
// @version      0.1
// @description  try to help me build jenkins faster
// @author       Yeqown@gmial.com
// @match        https://jenkins.offline-ops.net/*
// @grant        none
// @icon         https://www.jenkins.io/images/logos/jenkins/jenkins.svg
// @run-at       document-end
// ==/UserScript==

// 当前页面加载完成后，从中定位到 id="ecp_DeployMicroServices" 的元素
// 在前面插入两个按钮，一个是全选（id="selectAllButton"），一个是取消全选（id="deselectAllButton"）
// 全选的作用是，将下面这个列表中的所有的 checkbox 都选中
// 取消全选的作用是，将下面这个列表中的所有的 checkbox 都取消选中
// <table id="tbl_ecp_DeployMicroServices">
//  <tbody>
//   <tr style="white-space:nowrap" id="ecp_DeployMicroServices_0">
//    <td><input name="DeployMicroServices.value" json="game-api" checked="true" type="checkbox"
//            value="game-api"><label class="attach-previous">game-api</label></td>
//   </tr>
//   <tr style="white-space:nowrap" id="ecp_DeployMicroServices_1">
//    <td><input name="DeployMicroServices.value" json="game-cronjob" checked="true" type="checkbox"
//            value="game-cronjob"><label class="attach-previous">game-cronjob</label></td>
//   </tr>
//  </tbody>
// </table>

(function () {
    'use strict';

    console.log('Jenkins Hack running...');

    const table = document.getElementById('tbl_ecp_DeployMicroServices');
    if (!table) {
        console.warn('table(tbl_ecp_DeployMicroServices) not found');
        return;
    }

    const selectAllButton = document.createElement('button');
    selectAllButton.innerText = '全选';
    selectAllButton.id = 'selectAllButton';
    selectAllButton.onclick = function (evt) {
        evt.preventDefault();
        const checkboxes = table.querySelectorAll('input[type="checkbox"]');
        checkboxes.forEach(checkbox => {
            checkbox.checked = true;
        });
    };

    const deselectAllButton = document.createElement('button');
    deselectAllButton.innerText = '取消全选';
    deselectAllButton.id = 'deselectAllButton';
    deselectAllButton.onclick = function (evt) {
        evt.preventDefault();
        const checkboxes = table.querySelectorAll('input[type="checkbox"]');
        checkboxes.forEach(checkbox => {
            checkbox.checked = false;
        });
    };

    table.parentNode.insertBefore(selectAllButton, table);
    table.parentNode.insertBefore(deselectAllButton, table);
})();
