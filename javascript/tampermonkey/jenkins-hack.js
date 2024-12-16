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
// 再新增一个元素选中数量的说明：id="selectedCount"，内容为当前选中的数量 / 总数量
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
    const tableParent = document.getElementById('ecp_DeployMicroServices');
    if (!tableParent) {
        console.warn('table parent(ecp_DeployMicroServices) not found');
        return;
    }
    tableParent.style = 'float: left; padding-right: 25px;';

    const selectAllButton = document.createElement('button');
    selectAllButton.innerText = '全选';
    selectAllButton.id = 'selectAllButton';
    // selectAllButton.classList.add('jenkins-button');
    // selectAllButton.classList.add('jenkins-button--secondary');
    // selectAllButton.style.paddingRight = '1em';
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
    // deselectAllButton.classList.add('jenkins-button');
    // deselectAllButton.classList.add('jenkins-button--secondary');
    // deselectAllButton.style.paddingRight = '1em';
    deselectAllButton.onclick = function (evt) {
        evt.preventDefault();
        const checkboxes = table.querySelectorAll('input[type="checkbox"]');
        checkboxes.forEach(checkbox => {
            checkbox.checked = false;
        });
    };

    const selectedCount = document.createElement('span');
    selectedCount.id = 'selectedCount';
    selectedCount.innerText = '0 / 0';
    table.addEventListener('change', function (evt) {
        const checkboxes = table.querySelectorAll('input[type="checkbox"]');
        const selected = Array.from(checkboxes).filter(checkbox => checkbox.checked);
        selectedCount.innerText = `${selected.length} / ${checkboxes.length}`;
    });

    table.parentNode.insertBefore(selectAllButton, table);
    table.parentNode.insertBefore(deselectAllButton, table);
    table.parentNode.insertBefore(selectedCount, table);
})();
