'use strict';

const ROOT = 0;
const PARAGRAPH = 1;
const PHRASE = 2;
const WORD = 3;
const NUMBER = 4;
const SPACE = 5;
const OTHER = 6;


window.Leglog.notes = {};

class Text {
    id;
    title;
    root;
    language;
    creationTimestamp;

    constructor(id, creationTimestamp, language, title) {
        this.id = id;
        this.creationTimestamp = creationTimestamp;
        this.language = language;
        this.title = title;
    }

    initHTML() {
        this.root.initTag();
    }

    static createPhrase(nodes, phraseNode, note) {
        phraseNode.initTag();
        phraseNode.setNote(note);
        phraseNode.prev = nodes[0].prev;
        if (nodes[0].prev) {
            nodes[0].prev.next = phraseNode;
        }
        nodes[0].prev = null;
        phraseNode.next = nodes[nodes.length - 1].next;
        if (phraseNode.next) {
            phraseNode.next.prev = phraseNode;
        }
        nodes[nodes.length - 1].next = null
        phraseNode.children = nodes;
        phraseNode.parent = nodes[0].parent;
        phraseNode.tag.insertBefore(nodes[0].tag);
        for (const node of nodes) {
            node.parent = phraseNode;
            node.tag.appendTo(phraseNode.tag);
        }
    }
}

class Node {
    id;
    _text;
    prev;
    prevId;
    next;
    children;
    parent;
    parentId;
    noteId;
    note;
    type;
    tag;

 

    constructor(id, type, parentId, prevId, noteId, text) {
        this.id = id;
        this.type = type;
        this.parentId = parentId;
        this.prevId = prevId;
        this.noteId = noteId;
        this._text = text;
        this.children = [];
    }

    get isSpace() {
        return this.type == SPACE;
    }

    loadNote() {
        console.log(this.noteId);
        $.ajax('/note/load', {
            async: false,
            data: {
                noteId: this.noteId,
            },
            success: (data) => {
                console.log(data);
                let note = new Note(data['Id'], data['EntryId'], data['CreationTimestamp'],
                    data['UpdateTimestamp'], data['Content'], data['EntryText']);
                if (note.id in window.Leglog.notes) {
                    window.Leglog.notes[note.id].update(note);
                } else {
                    window.Leglog.notes[note.id] = note;
                }
                this.note = window.Leglog.notes[note.id];
            }
        });
    }

    setNote(note) {
        this.noteId = note.id;
        this.note = note;
    }

    get text() {
        if (this.children.length == 0) {
            return this._text;
        }
        let t = '';
        for (const child of this.children) {
            t += child.text;
        }
        return t;
    }
    //TODO: wrong!
    justPrecedes(other) {
        if(other==null){
            return false;
        }
        // if (this.type!=other.type)return false;
        let node = this.next;
        while (node != null && node != other && node.isSpace) {
            node = node.next;
        }
        return node == other;
    }

    justFollows(other) {
        if(other==null){
            return false;
        }
        // if (this.type!=other.type)return false;
        let node = this.prev;
        while (node != null && node != other && node.isSpace) {
            node = node.prev;
        }
        return node == other;
    }

    sortChildren() {
        children = [];
        let child;
        for (child of this.children) {
            if (child.prev === null) {
                break;
            }
        }
        while (child) {
            children.push(child);
            child = child.next;
        }
        this.children = children;
    }

    static Load() {

    }

    initTag() {
        switch (this.type) {
            case ROOT:
                this.tag = $('#divText');
                for (const child of this.children) {
                    child.initTag();
                    this.tag.append(child.tag);
                }
                this.tag.data('node', this);
                break;
            case PARAGRAPH:
                this.tag = $('<p></p>');
                for (const child of this.children) {
                    child.initTag();
                    this.tag.append(child.tag);
                }
                this.tag.data('node', this);
                break;
            case PHRASE:
                this.tag = $('<span></span>');
                this.tag.addClass('phrase');
                for (const child of this.children) {
                    child.initTag();
                    this.tag.append(child.tag);
                }
                this.tag.data('node', this);
                break;
            case WORD:
                this.tag = $('<span></span>');
                this.tag.addClass('word');
                this.tag.text(this._text);
                this.tag.data('node', this);
                break;
            default:
                this.tag = $(document.createTextNode(this._text));
                break;
        }
    }

}

class Entry {
    id;
    text;
    language;
    creationTimestamp;
    notes;

    loadNotes() {

    }
}

class Note {
    id;
    entryId;
    content;
    creationTimestamp;
    updateTimestamp;
    entryText;


    static fromObject(o){
        return new Note(
            o['Id'],
            o['EntryId'],
            o['CreationTimestamp'],
            o['UpdateTimestamp'],
            o['Content'],
            o['EntryText']
        );
    }

    constructor(id, entryId, creationTimestamp, updateTimestamp, content, entryText = null) {
        this.id = id;
        this.entryId = entryId;
        this.creationTimestamp = creationTimestamp;
        this.updateTimestamp = updateTimestamp;
        this.content = content;
        this.entryText = entryText
    }

    update(other) {
        this.id = other.id;
        this.entryId = other.entryId;
        this.content = other.content;
        this.creationTimestamp = other.creationTimestamp;
        this.updateTimestamp = other.updateTimestamp;
        this.entryText = other.entryText;
    }

    static Load() {

    }
}