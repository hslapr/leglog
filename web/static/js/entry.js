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

});