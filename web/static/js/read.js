'use strict';

function toggleOffcanvasLeft() {
    $('#offcanvasLeft').offcanvas('toggle');
}

function getLemmas() {
    $.ajax('/read/get-lemmas', {
        async: false,
        data: { word: window.Leglog.selection.text, language: window.Leglog.text.language },
        success: function (data) {
            console.log(data);
            if (data && data.length > 0) {
                for (const lemma of data) {
                    let option = $('<option></option>');
                    option.addClass('dropdown-item').attr('value', lemma).text(lemma);
                    $('#lemmaSelect').append(option);
                }
            }
        }
    });
}

function setCommentOptions() {
    let nComments = [
        'pl',
    ];
    let aComments = [
        'm', 'f', 'pl',
        'm.pl', 'f.pl',
    ];
    let vComments = [
        '1.sg.i.present', '2.sg.i.present', '3.sg.i.present',
        '1.pl.i.present', '2.pl.i.present', '3.pl.i.present',
        '1.sg.i.imperfect', '2.sg.i.imperfect', '3.sg.i.imperfect',
        '1.pl.i.imperfect', '2.pl.i.imperfect', '3.pl.i.imperfect',
        '1.sg.i.future', '2.sg.i.future', '3.sg.i.future',
        '1.pl.i.future', '2.pl.i.future', '3.pl.i.future',
    ];

    let comments = aComments.concat(vComments);
    for (const comment of comments) {
        let option = $('<option></option>');
        option.addClass('dropdown-item').attr('value', comment).text(comment);
        $('#commentSelect').append(option);
    }

}

function displayNote(note) {
    $('#divNotes').hide();
    $('#newNoteForm').hide();
    console.log(note);
    $('#editNoteTitle').text(note.entryText);
    $('#taEditNote').val(note.content);
    $('#btnSaveChange').hide();
    $('#divEditNote').show();
}

function queryNotes(s, lang) {
    $('#divEditNote').hide();
    $('#noteList').empty();
    $.ajax('/note/query', {
        data: {
            text: s,
            language: lang,
        },
        success: function (data) {
            console.log(data);
            if (data && data.length > 0) {
                for (const d of data) {
                    let note = new Note(d['Id'], d['EntryId'], d['ParentId'], d['PrevId'], d['Content'], d['EntryText']);
                    let a = $('<a></a>');
                    a.data('note', note);
                    a.addClass('list-group-item list-group-item-action note').attr('href', '#');
                    let div = $('<div></div>');
                    div.text('[' + note.entryText + '] ' + note.content);
                    a.append(div);
                    $('#noteList').append(a);
                }
                $('#newNoteForm').hide();
                $('#divNotes').show();
            } else {
                showNewNoteForm();
            }
        },
    });
}

function showNewNoteForm() {
    $('#divNotes').hide();
    $('#divEditNote').hide();
    $('#iLemma').val('');
    $('#taNewNote').val('');
    $('#lemmaSelect').empty();
    getLemmas();
    setCommentOptions();
    $('#newNoteForm').show();
}

class Selection {
    start;
    end;

    get nodes() {
        let nodes = [];
        let node = this.start;
        while (true) {
            nodes.push(node);
            if (node == this.end) {
                break;
            }
            node = node.next;
        }
        return nodes;
    }

    get isEmpty() {
        return this.start == null && this.end == null;
    }

    get isSingle() {
        return this.start != null && this.start == this.end;
    }

    get isMultiple() {
        return !this.isEmpty && !this.isSingle;
    }

    get text() {
        let t = '';
        let node = this.start;
        while (true) {
            t += node.text;
            if (node == this.end) {
                break;
            }
            node = node.next;
        }
        return t;
    }

    refresh() {
        if (this.start.prev == null && this.end.next == null) {
            this.selectSingle(this.start.parent);
        }
        if (this.isSingle && this.start.noteId > 0) {
            console.log(this.start.note);
            if (this.start.note == null) {
                this.start.loadNote();
            }
            displayNote(this.start.note);
        } else {
            queryNotes(this.text, window.Leglog.text.language);
        }

        // 方便复制
        $('#iSelectedText').val(this.text);
        $('#iSelectedText')[0].select();
        $('#iSelectedText').focus();
    }

    select(node) {
        if (node.justFollows(this.end)) {
            this.end = node;
            node.tag.addClass('selected');
        } else if (node.justPrecedes(this.start)) {
            this.start = node;
            node.tag.addClass('selected');
        } else {
            let parent = node.parent;
            while (parent != null) {
                if (parent.justPrecedes(this.start) || parent.justFollows(this.end)) {
                    this.select(parent);
                    break;
                }
                parent = parent.parent;
            }
            if (parent == null) {
                this.selectSingle(node.parent);
            }
        }
    }

    selectSingle(node) {
        this.clear();
        if (node == null) {
            return;
        }
        this.start = node;
        this.end = node;
        node.tag.addClass('selected');
    }

    clear() {
        if (this.isEmpty) {
            return;
        }
        let node = this.start;
        while (true) {
            node.tag.removeClass('selected');
            if (node == this.end) {
                break;
            }
            node = node.next;
        }
        this.start = null;
        this.end = null;
    }
}

function deleteText() {
    $('<form action="/text/delete" method="POST">' +
        '<input type="hidden" name="textId" value="' + window.Leglog.text.id + '">' +
        '<input type="hidden" name="rootId" value="' + window.Leglog.text.root.id + '">' +
        '</form>').appendTo('body').submit();
}


var selection;

$(() => {
    selection = new Selection();
    window.Leglog.selection = selection;
    window.Leglog.text.initHTML();


    // $('body').click(function(e){
    //     window.Leglog.selection.clear();
    // });

    $('body').on('click', 'a.note', function (e) {
        e.preventDefault();
        let note = $(this).data('note');
        let data = {
            Nodes: [],
            Note: note,
        };
        for (const node of window.Leglog.selection.nodes) {
            data.Nodes.push({
                Id: node.id,
                ParentId: node.parentId,
                PrevId: node.prevId,
                Text: node._text,
                NoteId: node.noteId,
                NodeType: node.type,
            });
        }
        console.log(data);
        $.ajax('/note/bind', {
            async: false,
            data: JSON.stringify(data),
            method: 'POST',
            contentType: "application/json",
            success: function (data) {
                console.log(data);
                let note = Note.fromObject(data['note']);
                if (note.id in window.Leglog.notes) {
                    window.Leglog.notes[note.id].update(note);
                } else {
                    window.Leglog.notes[note.id] = note;
                }
                if (window.Leglog.selection.isSingle) {
                    window.Leglog.selection.start.setNote(window.Leglog.notes[note.id]);
                } else {
                    let phraseNode = new Node(data['phraseNode']['Id'],
                        data['phraseNode']['NodeType'],
                        data['phraseNode']['ParentId'],
                        data['phraseNode']['PrevId'],
                        data['phraseNode']['NoteId'],
                        data['phraseNode']['Text']);
                    Text.createPhrase(window.Leglog.selection.nodes, phraseNode, window.Leglog.notes[note.id]);
                }
                window.Leglog.selection.refresh();
            },
        });
    });

    $('body').on('click', 'span.word', function (e) {
        if (e.ctrlKey) {
            window.Leglog.selection.select($(this).data('node'));
        } else {
            window.Leglog.selection.selectSingle($(this).data('node'));
        }
        window.Leglog.selection.refresh();
    });

    $('#btnSubmitNewNote').click(function (e) {
        e.preventDefault();
        let nodes = window.Leglog.selection.nodes;
        let entryText = window.Leglog.selection.text;
        let entryCase = $('#entryCaseSelect').val();
        switch (entryCase) {
            case '1':
                entryText = entryText.toLowerCase();
                break;
            case '2':
                // TODO: title case
                break;
            case '3':
                entryText = entryText.toUpperCase();
                break;
            default:
                break;
        }
        let lemma = $('#iLemma').val();
        let comment = $('#iComment').val();
        if ($('#taNewNote').val().length < 1) {
            $.ajax('/entry/create', {
                data: {
                    text: entryText,
                    lemma: lemma,
                    language: window.Leglog.text.language,
                    comment: comment
                },
                method: 'POST',
                success: function (e) {
                    window.Leglog.selection.refresh();
                },
            });
            return;
        }
        let data = {
            Nodes: [],
            Content: $('#taNewNote').val(),
            EntryText: entryText,
            Language: window.Leglog.text.language,
            Lemma: lemma,
            Comment: comment
        }
        for (const node of nodes) {
            data.Nodes.push({
                Id: node.id,
                ParentId: node.parentId,
                PrevId: node.prevId,
                Text: node._text,
                NoteId: node.noteId,
                NodeType: node.type,
            });
        }
        console.log(data);
        $.ajax('/note/create', {
            data: JSON.stringify(data),
            method: 'POST',
            async: false,
            contentType: "application/json",
            // dataType:'json',
            success: function (data) {
                console.log(data)
                let note = new Note(
                    data['note']['Id'],
                    data['note']['EntryId'],
                    data['note']['CreationTimestamp'],
                    data['note']['UpdateTimestamp'],
                    data['note']['Content'],
                    data['note']['EntryText']
                );
                if (note.id in window.Leglog.notes) {
                    window.Leglog.notes[note.id].update(note);
                } else {
                    window.Leglog.notes[note.id] = note;
                }
                if (window.Leglog.selection.isSingle) {
                    window.Leglog.selection.start.setNote(window.Leglog.notes[note.id]);
                } else {
                    let phraseNode = new Node(data['phraseNode']['Id'],
                        data['phraseNode']['NodeType'],
                        data['phraseNode']['ParentId'],
                        data['phraseNode']['PrevId'],
                        data['phraseNode']['NoteId'],
                        data['phraseNode']['Text']);
                    Text.createPhrase(window.Leglog.selection.nodes, phraseNode, window.Leglog.notes[note.id]);
                }
                window.Leglog.selection.refresh();
            },
        });
    });



    $('#btnNewNote').click(function (e) {
        $('#divNotes').hide();
        showNewNoteForm();
    });

    $('#taEditNote').change(function (e) {
        $('#btnSaveChange').show();
    });

    $('#btnSaveChange').click(function (e) {
        e.preventDefault();
        window.Leglog.selection.start.note.content = $('#taEditNote').val();
        $.ajax('/note/update', {
            data: JSON.stringify(window.Leglog.selection.start.note),
            method: 'POST',
            contentType: "application/json",
            success: function (data) {
                $('#btnSaveChange').hide();
                window.Leglog.selection.start.note.updateTimestamp = data['UpdateTimestamp'];
            }
        })
    });

    $('#btnRemoveNote').click(function (e) {
        e.preventDefault();
        let node = window.Leglog.selection.start;
        node.noteId = 0;
        node.note = null;
        $.ajax('/note/unbind', {
            data: JSON.stringify({
                id: node.id,
                nodeType: node.type,
                parentId: node.parentId,
                prevId: node.prevId,
                noteId: node.noteId,
            }),
            method: 'POST',
            contentType: "application/json",
            success: function (data) {
                if (node.type == PHRASE) {
                    for (let i = 0; i < node.children.length; i++) {
                        let child = node.children[i];
                        if (i == 0) {
                            child.prev = node.prev;
                            child.prevId = node.prevId;
                            if (child.prev) {
                                child.prev.next = child;
                            }
                            selection.start = child;
                        }
                        if (i == node.children.length - 1) {
                            child.next = node.next;
                            if (child.next) {
                                child.next.prev = child;
                                child.next.prevId = child.id;
                            }
                            selection.end = child;
                        }
                        child.parent = node.parent;
                        child.parentId = node.parentId;
                        if (child.parent) {
                            child.parent.children.splice(child.parent.children.indexOf(node), 0, child);
                        }
                        node.tag.before(child.tag);
                    }
                    if (node.parent) {
                        node.parent.children.splice(node.parent.children.indexOf(node, 1));
                    }
                    node.tag.remove();
                }
                window.Leglog.selection.refresh();
            }
        });
    });

    $('#btnConfirmLemma').click(function (e) {
        e.preventDefault();
        $('#iLemma').val($('#lemmaSelect').val());
    });

    $('#lemmaSelect').change(function (e) {
        $('#iLemma').val($('#lemmaSelect').val());
        $("#divLemmaDropdown").collapse('toggle');
    });

    $('#btnConfirmComment').click(function (e) {
        e.preventDefault();
        $('#iComment').val($('#commentSelect').val());
    });

    $('#commentSelect').change(function (e) {
        $('#iComment').val($('#commentSelect').val());
        $("#divCommentDropdown").collapse('toggle');
    });

    $('#btnDeleteText').click(function (e) {
        deleteText();
    });
});