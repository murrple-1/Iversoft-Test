'use strict';

function loadUsers() {
    $.ajax({
        type: 'GET',
        url: '/users',
        dataType: 'json'
    }).done(function (data) {
        var $table = $('table#users');
        var $tbody = $table.find('tbody');

        $tbody.empty();

        $.each(data, function (index, value) {
            var html = '<tr data-id="' + value.id + '">' +
                '<td>' + value.username + '</td>' +
                '<td>' + value.email + '</td>' +
                '<td>' +
                '<a href="#" class="edit-action"><i class="fas fa-edit"></i></a>' +
                '&nbsp;' +
                '<a href="#" class="delete-action"><i class="fas fa-trash"></i></a>' +
                '</td>' +
                '</tr>';

            $tbody.append(html);

            var $tr = $tbody.find('tr[data-id="' + value.id + '"]');

            $tr.find('a.delete-action').click(function () {
                deleteUser(value.id);
                return false;
            });

            $tr.find('a.edit-action').click(function () {
                showEditDialog(value.id);
                return false;
            })
        });
    }).fail(function () {
        console.log('error');
    });
}

function createUserHandler() {
    var $modal = $('div#create-user-modal');

    var payload = {};
    payload.username = $modal.find('#create-username-input').val();
    payload.email = $modal.find('#create-email-input').val();
    payload.roleLabel = $modal.find('#create-role-select').val();
    payload.address = {};
    payload.address.address = $modal.find('#create-address-input').val();
    payload.address.city = $modal.find('#create-city-input').val();
    payload.address.province = $modal.find('#create-province-input').val();
    payload.address.country = $modal.find('#create-country-input').val();
    payload.address.postalCode = $modal.find('#create-postal-code-input').val();

    $.ajax({
        type: 'POST',
        url: '/user',
        data: JSON.stringify(payload),
        contentType: 'application/json'
    }).done(function () {
        loadUsers();
        $modal.modal('hide');
    }).fail(function () {
        bootbox.alert('Failed to create user');
    });
}

function editUserHandler() {
    var $modal = $('div#edit-user-modal');

    var userId = $modal.data('userid');

    var payload = {};
    payload.email = $modal.find('#edit-email-input').val();
    payload.roleLabel = $modal.find('#edit-role-select').val();
    payload.address = {};
    payload.address.address = $modal.find('#edit-address-input').val();
    payload.address.city = $modal.find('#edit-city-input').val();
    payload.address.province = $modal.find('#edit-province-input').val();
    payload.address.country = $modal.find('#edit-country-input').val();
    payload.address.postalCode = $modal.find('#edit-postal-code-input').val();

    $.ajax({
        type: 'PUT',
        url: ('/user/' + userId),
        data: JSON.stringify(payload),
        contentType: 'application/json'
    }).done(function () {
        loadUsers();
        $modal.modal('hide');
    }).fail(function () {
        bootbox.alert('Failed to save user');
    });
}

function deleteUser(userId) {
    $.ajax({
        type: 'DELETE',
        url: ('/user/' + userId)
    }).done(function () {
        $('table#users tbody tr[data-id="' + userId + '"]').remove();
    }).fail(function () {
        bootbox.alert('Delete failed');
    });
}

function showEditDialog(userId) {
    $.ajax({
        type: 'GET',
        url: ('/user/' + userId),
        dataType: 'json'
    }).done(function (data) {
        var $modal = $('div#edit-user-modal');

        $modal.data('userid', userId);

        $modal.find('#edit-username-input').val(data.username);
        $modal.find('#edit-email-input').val(data.email);
        $modal.find('#edit-role-select').val(data.role.label);
        $modal.find('#edit-address-input').val(data.address.address);
        $modal.find('#edit-city-input').val(data.address.city);
        $modal.find('#edit-province-input').val(data.address.province);
        $modal.find('#edit-country-input').val(data.address.country);
        $modal.find('#edit-postal-code-input').val(data.address.postalCode);

        $modal.modal();
    }).fail(function () {
        bootbox.alert('Unable to retrieve user data');
    });
}

function clearCreateDialog() {
    var $modal = $('div#create-user-modal');

    $modal.find('#create-username-input').val('');
    $modal.find('#create-email-input').val('');
    $modal.find('#create-role-select').val('Public User');
    $modal.find('#create-address-input').val('');
    $modal.find('#create-city-input').val('');
    $modal.find('#create-province-input').val('');
    $modal.find('#create-country-input').val('');
    $modal.find('#create-postal-code-input').val('');
}

$(document).ready(function () {
    loadUsers();

    $('div.modal#create-user-modal').on('shown.bs.modal', clearCreateDialog);

    $('button#refresh-button').click(loadUsers);

    $('button#create-button').click(createUserHandler);

    $('button#edit-button').click(editUserHandler);
});
