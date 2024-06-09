// jQuery(function($) {
//   var $bodyEl = $('body'),
//       $sidedrawerEl = $('#sidedrawer');
  
  
//   // ==========================================================================
//   // Toggle Sidedrawer
//   // ==========================================================================
//   function showSidedrawer() {
//     // show overlay
//     var options = {
//       onclose: function() {
//         $sidedrawerEl
//           .removeClass('active')
//           .appendTo(document.body);
//       }
//     };
    
//     var $overlayEl = $(mui.overlay('on', options));
    
//     // show element
//     $sidedrawerEl.appendTo($overlayEl);
//     setTimeout(function() {
//       $sidedrawerEl.addClass('active');
//     }, 20);
//   }
  
  
//   function hideSidedrawer() {
//     $bodyEl.toggleClass('hide-sidedrawer');
//   }
  
  
//   $('.js-show-sidedrawer').on('click', showSidedrawer);
//   $('.js-hide-sidedrawer').on('click', hideSidedrawer);
  
  
//   // ==========================================================================
//   // Animate menu
//   // ==========================================================================
//   var $titleEls = $('strong', $sidedrawerEl);
  
//   $titleEls
//     .next()
//     .hide();
  
//   $titleEls.on('click', function() {
//     $(this).next().slideToggle(200);
//   });
// });
function exit(){
    Swal.fire({
      title: `Вы уверены что хотите выйти?`,
      icon: 'question',
      showCancelButton: true,
      confirmButtonText: 'Выйти',
      cancelButtonText: 'Отмена'
    })
    .then((result) => {
      if(result.isConfirmed){
        fetch("https://localost:443/private/sessionquit",{
          method:"GET",
          credentials:"include"
        })
        .then(response =>{
          if(response.ok){
          }
          else{
            console.log('failed to delete the session');
          }
        })
      }
      else{
        console.log('Выход отменён');
      }
    })
  }