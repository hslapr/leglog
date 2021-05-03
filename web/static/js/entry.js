'use strict'

$(()=>{
    $('button.btnDeleteNote').click(function(e){
        $.ajax('/note/delete', {
            data: {id: $(this).attr('data-note-id') },
            method: 'POST',
            success: function (data) {
                location.reload();
            }
        });
    });

    $('button.btnUpdateNote').click(function(e){
        $.ajax('/entry/update-note', {
            data: {id: $(this).attr('data-note-id'),content:$('#taNoteContent').val() },
            method: 'POST',
            success: function (data) {
                
            }
        });
    });

    $('button.btnRemoveLemma').click(function(e){
        $.ajax('/entry/remove-lemma', {
            data: {lemmaId: $(this).attr('data-lemma-id'), entryId: window.Leglog.entryId },
            method: 'POST',
            success: function (data) {
                location.reload();
            }
        });
    });

    $('#btnDeleteEntry').click(function(e){
        console.log('delete entry')
        $('<form action="/entry/delete" method="POST">' + 
              '<input type="hidden" name="id" value="' + window.Leglog.entryId + '">' +
              '</form>').appendTo('body').submit();
    });
});