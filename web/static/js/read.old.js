'use strict';
var selection;

function clearNewNoteForm() {
    $('#taNewNote').val('');
    $('iRoot').val('');

}

function showEditNote(note) {
    $('#iAddLemma').val('');
    $('#taEditNote').val(note['Content']);
    if (note['Root'] === null) {
        $('#editNoteTitle').text(note['Entry']['Name']);
        $('#btnSaveChange').removeClass('ms-auto').parent().removeClass('d-flex').addClass('input-group');
        $('#iAddLemma').show();
    } else {
        $('#editNoteTitle').text(note['Entry']['Name'] + ' (' + note['Root']['Name'] + ')');
        $('#iAddLemma').hide();
        $('#btnSaveChange').addClass('ms-auto').parent().removeClass('input-group').addClass('d-flex');
    }
    $('#divEditNote').attr('note-id', note['Id']);
    $('#divEditNote').show();
}

function loadEntries(s) {
    $.get('/read/get-entries', { word: s, language: lang }, function (data) {
        console.log(data)
        if (data === null || data.length == 0) {
            showNewNoteForm();
        } else {
            $('#divEntries button.btn-entry').remove();
            let btnNewEntryNote = $('#btnNewEntryNote');
            btnNewEntryNote.detach();
            for (const entry of data) {
                let btn = $('<button></button>');
                btn.addClass('btn-entry btn btn-success mx-1').text(entry['Name'])
                    .attr('entry-id', entry['Id']);
                $('#divEntries').append(btn);
            }
            $('#divEntries').append(btnNewEntryNote);
            $('#divEntries').show();
        }
    })
}

function loadNote(t) {
    if (t.data('note') === undefined) {
        $.ajax('/note/load', {
            data: { 'note-id': t.attr('note-id') },
            async: false,
            success: function (data) {
                t.data('note', data);
            }
        });
    }
}

function getLemmas() {
    $.ajax('/read/get-lemmas', {
        async: false,
        data: { word: selection.text, language: lang },
        success: function (data) {
            if (data != null && data.length > 0) {
                for (const lemma of data) {
                    let option = $('<option></option>');
                    option.addClass('dropdown-item').attr('value', lemma).text(lemma);
                    $('#lemmaSelect').append(option);
                }
            }
        }
    });
}

function showNewNoteForm() {
    $('#lemmaSelect').empty();
    $('#iRoot').val("");
    $('#taNewNote').val('');
    $('#entryCaseSelect').val(0);
    getLemmas();
    if ($('#newNoteForm').attr('entry-id') != undefined) {
        $('#entryCaseSelect').hide();
    } else {
        $('#entryCaseSelect').show();
    }
    $('#newNoteForm').show();
}

function loadNotes(entryId) {
    $.get('/note/get', { 'entry-id': entryId }, function (data) {
        if (data === null || data.length == 0) {
            showNewNoteForm();
        } else {
            console.log(data);
            $('#divNotes').show();
            $('#divNotes').attr('entry-id', entryId);
            $('#noteList').empty();
            for (const note of data) {
                let a = $('<a></a>');
                a.data('note', note);
                a.addClass('list-group-item list-group-item-action note').attr('href', '#')
                    .attr('note-id', note['Id']);
                let div = $('<div></div>');
                if (note['Root'] === null) {
                    div.text(note['Content']);
                } else {
                    div.text('[' + note['Root']['Name'] + '] ' + note['Content']);
                }
                a.append(div);
                $('#noteList').append(a);
            }
        }
    });
}

class Selection {
    constructor(start, end) {
        this.start = start;
        this.end = end;
    }

    setNoteId(noteId) {
        if (this.isSingleWord) {
            this.start.attr('note-id', noteId);
        } else {
            this.phraseTag.attr('note-id', noteId);
        }
    }

    get phraseTag() {
        return this.start.parent('.phrase')
    }

    get segId() {
        if (this.isSingleWord) {
            return this.start.attr('seg-id');
        } else {
            return this.phraseTag.attr('seg-id');
        }
    }

    get note() {
        if (this.isSingleWord) {
            return this.start.data('note');
        } else {
            return this.phraseTag.data('note');
        }
    }

    get text() {
        if (this.isEmpty) {
            return '';
        }
        let t = this.start[0];
        let s = '';
        while (t != this.end[0].nextSibling) {
            s += t.textContent;
            t = t.nextSibling;
        }
        return s;
    }


    change() {
        // 方便复制
        $('#iSelectedText').val(this.text);
        $('#iSelectedText')[0].select();
        $('#iSelectedText').focus();
        // 下面部分全部隐藏
        $('#divEntries').hide();
        $('#divNotes').hide();
        $('#newNoteForm').hide();
        $('#divEditNote').hide();
        $('#btnSaveChange').hide();
        $('#newNoteForm').removeAttr('entry-id');
        if (this.isSingleWord) {
            if (this.start.attr('note-id') === undefined) {
                loadEntries(selection.text);
            } else {
                loadNote(this.start);
                showEditNote(this.start.data('note'));
            }
        } else {
            if (this.phraseTag.length > 0) {
                loadNote(this.phraseTag);
                showEditNote(this.phraseTag.data('note'));
            } else {
                loadEntries(selection.text);
            }
        }
    }

    get isSingleWord() {
        return !this.isEmpty && this.start[0] == this.end[0]
    }

    setStart(s) {
        this.start = s;
        s.addClass('selected');
        this.change()
    }

    setEnd(s) {
        this.end = s;
        s.addClass('selected');
        this.change()
    }

    get isEmpty() {
        return this.start === null && this.end === null;
    }

    clear() {
        console.log('clear');
        let t = this.start;
        while (t[0] != this.end[0]) {
            t.removeClass('selected');
            t = t.next();
        }
        t.removeClass('selected');
        this.start = null;
        this.end = null;
    }

    select(t) {
        this.start = t;
        this.end = t;
        t.addClass('selected')
        this.change()
    }
}

function isPreceding(t1, t2) {
    if (t1.nextSibling === null) {
        return false;
    }
    if (t1.nextSibling.nodeName == '#text' && /^\s*$/.test(t1.nextSibling.textContent)) {
        return t1.nextSibling.nextSibling == t2;
    }
    return t1.nextSibling == t2;
}

$(() => {


    selection = new Selection(null, null)

    $('#btnConfirmLemma').click(function (e) {
        e.preventDefault();
        $('#iRoot').val($('#lemmaSelect').val());
    });

    $('#lemmaSelect').change(function (e) {
        $('#iRoot').val($('#lemmaSelect').val());
        $("#newNoteForm div.dropdown-menu").collapse('toggle');
    });

    $('#btnSaveChange').click(function (e) {
        e.preventDefault();
        $.ajax('/note/update', {
            method: 'POST',
            data: {
                noteId: $('#divEditNote').attr('note-id'),
                content: $('#taEditNote').val(),
                lemma: $('#iAddLemma').val(),
                language:lang
            },
            success: function (data) {
                $('#btnSaveChange').hide();
                selection.note['Content'] = $('#taEditNote').val();
            }
        })
    });

    $('#btnNewEntryNote').click(function (e) {
        $('#divEntries').hide();
        showNewNoteForm();
    });
    $('#btnNewNote').click(function (e) {
        $('#divNotes').hide();
        showNewNoteForm();
    });

    $('body').on('click', 'button.btn-entry', function (e) {
        let entryId = $(this).attr('entry-id');
        $('#newNoteForm').attr('entry-id', entryId);
        loadNotes(entryId);
        $('#divEntries').hide();
    });

    $('span.word').click(function (e) {
        if (selection.isEmpty) {
            selection.select($(this));
        } else {
            if (e.ctrlKey) {
                if (isPreceding(this, selection.start[0])) {
                    selection.setStart($(this));
                } else if (isPreceding(selection.end[0], this)) {
                    selection.setEnd($(this));
                } else {
                    selection.clear()
                    selection.select($(this));
                }
            } else {
                selection.clear()
                selection.select($(this));
            }
        }
    });
    $('body').on('click', 'a.note', function (e) {
        e.preventDefault();
        $.ajax('/read/set-segment-note', {
            method: 'POST',
            async: false,
            data: {
                segId: selection.segId,
                noteId: $(this).attr('note-id')
            },
            success: function (data) {
                selection.setNoteId(data["note-id"]);
                selection.change();
            }
        });
        $('#divNotes').hide();
        // let note = $(this).data('note');
        // showEditNote(note);
    });



    $('#taEditNote').change(function (e) {
        $('#btnSaveChange').show();
    });

    $('#btnSubmitNewNote').click(function (e) {
        e.preventDefault();
        if (!selection.isSingleWord) {
            $.ajax('/read/group-segments', {
                data: {
                    startId: selection.start.attr('seg-id'),
                    endId: selection.end.attr('seg-id')
                },
                method: 'POST',
                async: false,
                success: function (data) {
                    let span = $('<span></span>');
                    span.addClass('phrase').attr('seg-id', data['seg-id']);
                    selection.start.before(span);
                    let t = selection.start[0];
                    let next = t.nextSibling;
                    let end = selection.end[0].nextSibling
                    while (t != end) {
                        span[0].appendChild(t);
                        t = next;
                        if (next === null) {
                            break;
                        }
                        next = next.nextSibling;
                    }

                }
            })
        }

        $.ajax('/note/create', {
            data: {
                segId: selection.segId,
                content: $('#taNewNote').val(),
                word: selection.text,
                case: $('#entryCaseSelect').val(),
                root: $('#iRoot').val(),
                language: lang
            },
            method: 'POST',
            success: function (data) {
                selection.setNoteId(data['Id']);
                clearNewNoteForm();
                selection.change();
            }
        });
    })

})